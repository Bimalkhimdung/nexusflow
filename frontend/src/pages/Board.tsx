import { useState, useEffect } from 'react';
import { DragDropContext, Droppable, Draggable, DropResult } from '@hello-pangea/dnd';
import { useAppStore } from '@/stores';
import { Issue } from '@/types';
import { Plus, MoreHorizontal } from 'lucide-react';

const COLUMNS = [
    { id: 'TODO', title: 'To Do', color: 'bg-slate-100 dark:bg-slate-800' },
    { id: 'IN_PROGRESS', title: 'In Progress', color: 'bg-blue-50 dark:bg-blue-900/20' },
    { id: 'DONE', title: 'Done', color: 'bg-green-50 dark:bg-green-900/20' }
];

export function Board() {
    const { issues, setIssues, fetchIssues, projects } = useAppStore();
    const [columns, setColumns] = useState<{ [key: string]: Issue[] }>({});

    useEffect(() => {
        if (projects.length > 0) {
            fetchIssues(projects[0].id);
        }
    }, [projects]);

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

        // TODO: Call API to update status
    };

    return (
        <div className="h-full flex flex-col">
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Board</h2>
                    <p className="text-muted-foreground">Manage tasks visually.</p>
                </div>
                <button className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
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
        </div>
    );
}
