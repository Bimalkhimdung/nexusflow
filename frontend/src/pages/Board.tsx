import { useState, useEffect } from 'react';
import { DragDropContext, Droppable, Draggable, DropResult } from '@hello-pangea/dnd';
import { useAppStore } from '@/stores';
import { Issue } from '@/types';
import { Plus, MoreHorizontal } from 'lucide-react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';

const issueSchema = z.object({
    title: z.string().min(1, 'Title is required'),
    description: z.string().optional(),
    priority: z.enum(['LOW', 'MEDIUM', 'HIGH']),
    type: z.enum(['TASK', 'BUG', 'STORY', 'EPIC'])
});

const COLUMNS = [
    { id: 'TODO', title: 'To Do', color: 'bg-slate-100 dark:bg-slate-800' },
    { id: 'IN_PROGRESS', title: 'In Progress', color: 'bg-blue-50 dark:bg-blue-900/20' },
    { id: 'DONE', title: 'Done', color: 'bg-green-50 dark:bg-green-900/20' }
];


type IssueFormValues = z.infer<typeof issueSchema>;

export function Board() {
    const { issues, setIssues, fetchIssues, createIssue, projects, currentProject, fetchProjects, updateIssue } = useAppStore();
    const [columns, setColumns] = useState<{ [key: string]: Issue[] }>({});
    const [isModalOpen, setIsModalOpen] = useState(false);
    const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<IssueFormValues>({
        resolver: zodResolver(issueSchema),
        defaultValues: {
            priority: 'MEDIUM',
            type: 'TASK'
        }
    });

    // Fetch projects first
    useEffect(() => {
        fetchProjects();
    }, []);

    useEffect(() => {
        if (currentProject) {
            fetchIssues(currentProject.id);
        } else if (projects.length > 0) {
            fetchIssues(projects[0].id);
        }
    }, [projects, currentProject]);

    useEffect(() => {
        const newColumns: { [key: string]: Issue[] } = {
            TODO: [],
            IN_PROGRESS: [],
            DONE: []
        };

        issues.forEach(issue => {
            if (newColumns[issue.status]) {
                newColumns[issue.status].push(issue);
            }
        });

        setColumns(newColumns);
    }, [issues]);

    const onDragEnd = (result: DropResult) => {
        const { source, destination, draggableId } = result;

        if (!destination) return;

        if (
            source.droppableId === destination.droppableId &&
            source.index === destination.index
        ) {
            return;
        }

        const sourceColumn = columns[source.droppableId];
        const movedIssue = sourceColumn.find(i => i.id === draggableId);

        if (!movedIssue) return;

        // Optimistic update
        const newIssues = [...issues];
        const issueIndex = newIssues.findIndex(i => i.id === draggableId);
        if (issueIndex !== -1) {
            newIssues[issueIndex] = {
                ...newIssues[issueIndex],
                status: destination.droppableId as 'TODO' | 'IN_PROGRESS' | 'DONE'
            };
            setIssues(newIssues);
        }

        // Call API to update status
        updateIssue(draggableId, {
            status: destination.droppableId as 'TODO' | 'IN_PROGRESS' | 'DONE'
        });
    };

    const onSubmit = async (data: IssueFormValues) => {
        try {
            const projectId = currentProject?.id || projects[0]?.id;
            if (!projectId) {
                console.error('No project selected. Please create a project first.');
                alert('Please create a project first before creating issues.');
                return;
            }
            await createIssue({
                ...data,
                summary: data.title,
                type: `ISSUE_TYPE_${data.type}` as any,
                priority: `ISSUE_PRIORITY_${data.priority}` as any,
                projectId,
                status: 'TODO'
            });
            setIsModalOpen(false);
            reset();
        } catch (error) {
            console.error('Failed to create issue:', error);
            alert('Failed to create issue. Please check console for details.');
        }
    };

    return (
        <div className="h-full flex flex-col">
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Board</h2>
                    <p className="text-muted-foreground">Manage tasks visually.</p>
                </div>
                <button
                    onClick={() => setIsModalOpen(true)}
                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
                >
                    <Plus className="mr-2 h-4 w-4" />
                    New Issue
                </button>
            </div>

            <DragDropContext onDragEnd={onDragEnd}>
                <div className="flex h-full gap-6 overflow-x-auto pb-4">
                    {COLUMNS.map(column => (
                        <div key={column.id} className={`flex-shrink-0 w-80 flex flex-col rounded-lg border bg-card text-card-foreground shadow-sm`}>
                            <div className="p-4 font-semibold flex items-center justify-between border-b">
                                <div className="flex items-center gap-2">
                                    <div className={`w-3 h-3 rounded-full ${column.id === 'TODO' ? 'bg-slate-400' :
                                        column.id === 'IN_PROGRESS' ? 'bg-blue-500' : 'bg-green-500'
                                        }`} />
                                    {column.title}
                                    <span className="ml-2 rounded-full bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
                                        {columns[column.id]?.length || 0}
                                    </span>
                                </div>
                                <button className="text-muted-foreground hover:text-foreground">
                                    <MoreHorizontal className="h-4 w-4" />
                                </button>
                            </div>

                            <Droppable droppableId={column.id}>
                                {(provided, snapshot) => (
                                    <div
                                        {...provided.droppableProps}
                                        ref={provided.innerRef}
                                        className={`flex-1 p-2 space-y-3 transition-colors ${snapshot.isDraggingOver ? 'bg-muted/50' : ''
                                            }`}
                                    >
                                        {columns[column.id]?.map((issue, index) => (
                                            <Draggable key={issue.id} draggableId={issue.id} index={index}>
                                                {(provided, snapshot) => (
                                                    <div
                                                        ref={provided.innerRef}
                                                        {...provided.draggableProps}
                                                        {...provided.dragHandleProps}
                                                        className={`rounded-md border bg-background p-3 shadow-sm transition-shadow hover:shadow-md ${snapshot.isDragging ? 'shadow-lg ring-2 ring-primary ring-opacity-50 rotate-2' : ''
                                                            }`}
                                                        style={provided.draggableProps.style}
                                                    >
                                                        <div className="flex flex-col gap-2">
                                                            <div className="flex items-start justify-between">
                                                                <span className="font-medium text-sm hover:underline cursor-pointer">
                                                                    {issue.title}
                                                                </span>
                                                            </div>
                                                            <div className="flex items-center justify-between text-xs text-muted-foreground">
                                                                <div className="flex items-center gap-2">
                                                                    <span className="font-mono text-[10px] uppercase bg-muted px-1 rounded">
                                                                        {issue.key}
                                                                    </span>
                                                                    <div className={`w-2 h-2 rounded-full ${issue.priority === 'HIGH' ? 'bg-red-500' :
                                                                        issue.priority === 'MEDIUM' ? 'bg-yellow-500' : 'bg-blue-500'
                                                                        }`} title={`Priority: ${issue.priority}`} />
                                                                </div>
                                                                {issue.assigneeId && (
                                                                    <div className="w-5 h-5 rounded-full bg-primary/10 flex items-center justify-center text-[10px] font-bold text-primary">
                                                                        {issue.assigneeId.charAt(0)}
                                                                    </div>
                                                                )}
                                                            </div>
                                                        </div>
                                                    </div>
                                                )}
                                            </Draggable>
                                        ))}
                                        {provided.placeholder}
                                    </div>
                                )}
                            </Droppable>
                        </div>
                    ))}
                </div>
            </DragDropContext>

            {isModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
                    <div className="w-full max-w-md rounded-lg bg-background p-6 shadow-lg animate-in fade-in zoom-in duration-200">
                        <h3 className="text-lg font-semibold mb-4">Create New Issue</h3>
                        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Title</label>
                                <input
                                    {...register('title')}
                                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    placeholder="Issue title"
                                />
                                {errors.title && <p className="text-sm text-red-500">{errors.title.message}</p>}
                            </div>

                            <div className="space-y-2">
                                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Description</label>
                                <textarea
                                    {...register('description')}
                                    className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    placeholder="Issue description"
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Type</label>
                                    <select
                                        {...register('type')}
                                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    >
                                        <option value="TASK">Task</option>
                                        <option value="BUG">Bug</option>
                                        <option value="STORY">Story</option>
                                        <option value="EPIC">Epic</option>
                                    </select>
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Priority</label>
                                    <select
                                        {...register('priority')}
                                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    >
                                        <option value="LOW">Low</option>
                                        <option value="MEDIUM">Medium</option>
                                        <option value="HIGH">High</option>
                                    </select>
                                </div>
                            </div>

                            <div className="flex justify-end gap-2 mt-6">
                                <button
                                    type="button"
                                    onClick={() => {
                                        setIsModalOpen(false);
                                        reset();
                                    }}
                                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
                                >
                                    {isSubmitting ? 'Creating...' : 'Create Issue'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
