import { create } from 'zustand';
import { Issue, Project, User } from '../types';


interface AppState {
    user: User | null;
    projects: Project[];
    currentProject: Project | null;
    issues: Issue[];
    loading: boolean;

    setUser: (user: User | null) => void;
    setProjects: (projects: Project[]) => void;
    setCurrentProject: (project: Project | null) => void;
    setIssues: (issues: Issue[]) => void;

    fetchProjects: () => Promise<void>;
    fetchIssues: (projectId: string) => Promise<void>;
}

export const useAppStore = create<AppState>((set) => ({
    user: {
        id: '1',
        email: 'demo@nexusflow.io',
        fullName: 'Demo User',
        role: 'ADMIN'
    }, // Mock user for now
    projects: [],
    currentProject: null,
    issues: [],
    loading: false,

    setUser: (user) => set({ user }),
    setProjects: (projects) => set({ projects }),
    setCurrentProject: (currentProject) => set({ currentProject }),
    setIssues: (issues) => set({ issues }),

    fetchProjects: async () => {
        set({ loading: true });
        try {
            // Mock data for now until backend is fully connected via gateway
            const mockProjects: Project[] = [
                {
                    id: '1',
                    name: 'NexusFlow',
                    key: 'NEX',
                    description: 'Project Management System',
                    ownerId: '1',
                    createdAt: new Date().toISOString()
                }
            ];
            set({ projects: mockProjects, currentProject: mockProjects[0] });
        } catch (error) {
            console.error('Failed to fetch projects', error);
        } finally {
            set({ loading: false });
        }
    },

    fetchIssues: async (projectId) => {
        set({ loading: true });
        try {
            // Mock data
            const mockIssues: Issue[] = [
                {
                    id: '1',
                    key: 'NEX-1',
                    title: 'Implement Frontend',
                    description: 'Build the React frontend',
                    status: 'IN_PROGRESS',
                    priority: 'HIGH',
                    type: 'STORY',
                    reporterId: '1',
                    projectId: projectId,
                    createdAt: new Date().toISOString(),
                    updatedAt: new Date().toISOString()
                },
                {
                    id: '2',
                    key: 'NEX-2',
                    title: 'Design Database',
                    description: 'Create schema',
                    status: 'DONE',
                    priority: 'MEDIUM',
                    type: 'TASK',
                    reporterId: '1',
                    projectId: projectId,
                    createdAt: new Date().toISOString(),
                    updatedAt: new Date().toISOString()
                }
            ];
            set({ issues: mockIssues });
        } catch (error) {
            console.error('Failed to fetch issues', error);
        } finally {
            set({ loading: false });
        }
    }
}));
