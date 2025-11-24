import { create } from 'zustand';
import { Issue, Project, User } from '../types';
import api from '../services/api';

// Hardcoded organization ID for demo purposes (replace with dynamic value as needed)
const ORG_ID = 'e5616b48-3cdf-4812-a54f-fd861d1ff062';

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

    createProject: (project: Partial<Project>) => Promise<void>;
    updateProject: (projectId: string, project: Partial<Project>) => Promise<void>;
    deleteProject: (projectId: string) => Promise<void>;
    fetchProjects: () => Promise<void>;
    fetchIssues: (projectId: string) => Promise<void>;
    createIssue: (issue: Partial<Issue>) => Promise<void>;
}

export const useAppStore = create<AppState>((set, get) => ({
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

    createProject: async (project: Partial<Project>) => {
        set({ loading: true });
        try {
            const response = await api.post('/projects', { ...project, organization_id: ORG_ID });
            set((state) => ({ projects: [...state.projects, response.data.project] }));
        } catch (error) {
            console.error('Failed to create project', error);
            throw error;
        } finally {
            set({ loading: false });
        }
    },

    deleteProject: async (projectId: string) => {
        set({ loading: true });
        try {
            await api.delete(`/projects/${projectId}`);
            set((state) => ({
                projects: state.projects.filter(p => p.id !== projectId),
                currentProject: state.currentProject?.id === projectId ? null : state.currentProject
            }));
        } catch (error) {
            console.error('Failed to delete project', error);
            throw error;
        } finally {
            set({ loading: false });
        }
    },

    updateProject: async (projectId: string, project: Partial<Project>) => {
        set({ loading: true });
        try {
            const response = await api.patch(`/projects/${projectId}`, project);
            set((state) => ({
                projects: state.projects.map(p =>
                    p.id === projectId ? response.data.project : p
                ),
                currentProject: state.currentProject?.id === projectId ? response.data.project : state.currentProject
            }));
        } catch (error) {
            console.error('Failed to update project', error);
            throw error;
        } finally {
            set({ loading: false });
        }
    },

    fetchProjects: async () => {
        set({ loading: true });
        try {
            const response = await api.get(`/projects?organization_id=${ORG_ID}`);
            set({ projects: response.data.projects || [] });
            if (response.data.projects && response.data.projects.length > 0 && !get().currentProject) {
                set({ currentProject: response.data.projects[0] });
            }
        } catch (error) {
            console.error('Failed to fetch projects', error);
            set({ projects: [] }); // Ensure projects is always an array
        } finally {
            set({ loading: false });
        }
    },

    fetchIssues: async (projectId) => {
        set({ loading: true });
        try {
            const response = await api.get(`/projects/${projectId}/issues`);
            const issues = (response.data.issues || []).map((issue: any) => ({
                ...issue,
                title: issue.summary || 'Untitled',
                status: 'TODO', // Default to TODO as backend uses statusId
                priority: issue.priority ? issue.priority.replace('ISSUE_PRIORITY_', '') : 'MEDIUM',
                type: issue.type ? issue.type.replace('ISSUE_TYPE_', '') : 'TASK'
            }));
            set({ issues });
        } catch (error) {
            console.error('Failed to fetch issues', error);
            set({ issues: [] });
        } finally {
            set({ loading: false });
        }
    },

    createIssue: async (issue: Partial<Issue>) => {
        set({ loading: true });
        try {
            const response = await api.post(`/projects/${issue.projectId}/issues`, issue);
            const newIssue = {
                ...response.data.issue,
                title: response.data.issue.summary || issue.title || 'Untitled',
                status: 'TODO',
                priority: response.data.issue.priority ? response.data.issue.priority.replace('ISSUE_PRIORITY_', '') : 'MEDIUM',
                type: response.data.issue.type ? response.data.issue.type.replace('ISSUE_TYPE_', '') : 'TASK'
            };
            set((state) => ({ issues: [...state.issues, newIssue] }));
        } catch (error) {
            console.error('Failed to create issue', error);
            throw error;
        } finally {
            set({ loading: false });
        }
    }
}));
