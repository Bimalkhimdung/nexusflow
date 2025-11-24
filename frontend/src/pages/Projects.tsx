import { useState, useEffect } from 'react';
import { useAppStore } from '@/stores';
import { Plus, Folder, Calendar, Trash2, Pencil } from 'lucide-react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import type { Project } from '@/types';

const projectSchema = z.object({
    name: z.string().min(1, 'Name is required'),
    key: z.string().min(1, 'Key is required').max(5, 'Key must be 5 characters or less').toUpperCase(),
    description: z.string().optional(),
});

type ProjectFormValues = z.infer<typeof projectSchema>;

export function Projects() {
    const { projects, fetchProjects, createProject, updateProject, deleteProject, loading } = useAppStore();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isEditMode, setIsEditMode] = useState(false);
    const [editingProject, setEditingProject] = useState<Project | null>(null);
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
    const [projectToDelete, setProjectToDelete] = useState<string | null>(null);
    const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<ProjectFormValues>({
        resolver: zodResolver(projectSchema)
    });

    useEffect(() => {
        fetchProjects().catch(error => {
            console.error("Error fetching projects in component:", error);
        });
    }, [fetchProjects]);

    const onSubmit = async (data: ProjectFormValues) => {
        try {
            if (isEditMode && editingProject) {
                await updateProject(editingProject.id, data);
            } else {
                await createProject(data);
            }
            setIsModalOpen(false);
            setIsEditMode(false);
            setEditingProject(null);
            reset();
        } catch (error) {
            console.error(error);
        }
    };

    const handleEditClick = (project: Project) => {
        setEditingProject(project);
        setIsEditMode(true);
        reset({
            name: project.name,
            key: project.key,
            description: project.description || ''
        });
        setIsModalOpen(true);
    };

    const handleCreateClick = () => {
        setIsEditMode(false);
        setEditingProject(null);
        reset({ name: '', key: '', description: '' });
        setIsModalOpen(true);
    };

    const handleDeleteClick = (projectId: string) => {
        setProjectToDelete(projectId);
        setDeleteConfirmOpen(true);
    };

    const handleDeleteConfirm = async () => {
        if (projectToDelete) {
            try {
                await deleteProject(projectToDelete);
                setDeleteConfirmOpen(false);
                setProjectToDelete(null);
            } catch (error) {
                console.error('Failed to delete project:', error);
            }
        }
    };

    return (
        <div className="h-full flex flex-col space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Projects</h2>
                    <p className="text-muted-foreground">Manage your projects and teams.</p>
                </div>
                <button
                    onClick={handleCreateClick}
                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
                >
                    <Plus className="mr-2 h-4 w-4" />
                    New Project
                </button>
            </div>


            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {loading && <p>Loading projects...</p>}
                {!loading && (!projects || projects.length === 0) && (
                    <p className="text-muted-foreground">No projects found. Create your first project!</p>
                )}
                {!loading && projects && projects.map((project) => (
                    <div key={project.id} className="rounded-lg border bg-card text-card-foreground shadow-sm hover:shadow-md transition-shadow">
                        <div className="flex flex-col space-y-1.5 p-6">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Folder className="h-5 w-5 text-blue-500" />
                                    <h3 className="font-semibold leading-none tracking-tight">{project.name}</h3>
                                </div>
                                <span className="text-xs font-mono bg-muted px-2 py-1 rounded">{project.key}</span>
                            </div>
                            <p className="text-sm text-muted-foreground line-clamp-2 mt-2">
                                {project.description || 'No description provided.'}
                            </p>
                        </div>
                        <div className="p-6 pt-0 flex items-center justify-between text-sm text-muted-foreground">
                            <div className="flex items-center gap-1">
                                <Calendar className="h-4 w-4" />
                                <span>{new Date(project.createdAt).toLocaleDateString()}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <button
                                    onClick={() => handleEditClick(project)}
                                    className="text-blue-500 hover:text-blue-700 transition-colors p-1 rounded hover:bg-blue-50"
                                    title="Edit project"
                                >
                                    <Pencil className="h-4 w-4" />
                                </button>
                                <button
                                    onClick={() => handleDeleteClick(project.id)}
                                    className="text-red-500 hover:text-red-700 transition-colors p-1 rounded hover:bg-red-50"
                                    title="Delete project"
                                >
                                    <Trash2 className="h-4 w-4" />
                                </button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            {isModalOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
                    <div className="w-full max-w-md rounded-lg bg-background p-6 shadow-lg animate-in fade-in zoom-in duration-200">
                        <h3 className="text-lg font-semibold mb-4">{isEditMode ? 'Edit Project' : 'Create New Project'}</h3>
                        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Name</label>
                                <input
                                    {...register('name')}
                                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    placeholder="Project Name"
                                />
                                {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
                            </div>

                            <div className="space-y-2">
                                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Key</label>
                                <input
                                    {...register('key')}
                                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    placeholder="PRJ"
                                    maxLength={5}
                                />
                                {errors.key && <p className="text-sm text-red-500">{errors.key.message}</p>}
                            </div>

                            <div className="space-y-2">
                                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Description</label>
                                <textarea
                                    {...register('description')}
                                    className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                    placeholder="Project Description"
                                />
                            </div>

                            <div className="flex justify-end gap-2 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setIsModalOpen(false)}
                                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
                                >
                                    {isSubmitting ? 'Creating...' : 'Create Project'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {deleteConfirmOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
                    <div className="w-full max-w-md rounded-lg bg-background p-6 shadow-lg animate-in fade-in zoom-in duration-200">
                        <h3 className="text-lg font-semibold mb-4">Delete Project</h3>
                        <p className="text-muted-foreground mb-6">
                            Are you sure you want to delete this project? This action cannot be undone.
                        </p>
                        <div className="flex justify-end gap-2">
                            <button
                                onClick={() => {
                                    setDeleteConfirmOpen(false);
                                    setProjectToDelete(null);
                                }}
                                className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleDeleteConfirm}
                                disabled={loading}
                                className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-red-600 text-white hover:bg-red-700 h-10 px-4 py-2"
                            >
                                {loading ? 'Deleting...' : 'Delete'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
