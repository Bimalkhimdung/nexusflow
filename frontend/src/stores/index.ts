import { create } from 'zustand';
import { Issue, Project, User, Organization, OrgMember, Invite } from '../types';
import api from '../services/api';

// Hardcoded organization ID for demo purposes (replace with dynamic value as needed)
const ORG_ID = 'e5616b48-3cdf-4812-a54f-fd861d1ff062';

interface AppState {
    user: User | null;
    userRole: 'owner' | 'admin' | 'member' | 'guest' | null;
    organizations: Organization[];
    currentOrganization: Organization | null;
    orgMembers: OrgMember[];
    invites: Invite[];
    projects: Project[];
    currentProject: Project | null;
    issues: Issue[];
    loading: boolean;

    setUser: (user: User | null) => void;
    setProjects: (projects: Project[]) => void;
    setCurrentProject: (project: Project | null) => void;
    setIssues: (issues: Issue[]) => void;
    setCurrentOrganization: (org: Organization | null) => void;

    createProject: (project: Partial<Project>) => Promise<void>;
    updateProject: (projectId: string, project: Partial<Project>) => Promise<void>;
    deleteProject: (projectId: string) => Promise<void>;
    fetchProjects: () => Promise<void>;
    fetchIssues: (projectId: string) => Promise<void>;
    createIssue: (issue: Partial<Issue>) => Promise<void>;
    updateIssue: (issueId: string, issue: Partial<Issue>) => Promise<void>;

    // Organization & Member Management
    fetchUserRole: (orgId: string) => Promise<void>;
    fetchOrganizations: () => Promise<void>;
    fetchOrgMembers: (orgId: string) => Promise<void>;
    inviteMember: (orgId: string, email: string, role: 'admin' | 'member') => Promise<void>;
    removeMember: (orgId: string, userId: string) => Promise<void>;
    fetchInvites: (orgId: string) => Promise<void>;
    revokeInvite: (inviteId: string) => Promise<void>;
    addProjectMember: (projectId: string, userId: string, role: string) => Promise<void>;
    removeProjectMember: (projectId: string, userId: string) => Promise<void>;
}

export const useAppStore = create<AppState>((set, get) => ({
    user: {
        id: '1',
        email: 'demo@nexusflow.io',
        fullName: 'Demo User',
        role: 'ADMIN'
    }, // Mock user for now
    userRole: 'admin', // Mock role
    organizations: [],
    currentOrganization: null,
    orgMembers: [],
    invites: [],
    projects: [],
    currentProject: null,
    issues: [],
    loading: false,

    setUser: (user) => set({ user }),
    setProjects: (projects) => set({ projects }),
    setCurrentProject: (currentProject) => set({ currentProject }),
    setIssues: (issues) => set({ issues }),
    setCurrentOrganization: (currentOrganization) => set({ currentOrganization }),

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
                status: issue.statusId || 'TODO',
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
                status: response.data.issue.statusId || 'TODO',
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
    },

    updateIssue: async (issueId: string, issue: Partial<Issue>) => {
        // Optimistic update is handled by the component, but we should also update store here
        // to ensure consistency if the API call fails (we would revert) or succeeds
        set((state) => ({
            issues: state.issues.map((i) => (i.id === issueId ? { ...i, ...issue } : i))
        }));

        try {
            // Map frontend fields back to backend fields
            const payload: any = { ...issue };
            if (issue.title) payload.summary = issue.title;
            if (issue.status) payload.status_id = issue.status;
            // Status is tricky because backend uses statusId, but for now we might not have status mapping
            // If backend supports status update via specific endpoint or field, we use that.
            // Assuming backend accepts 'status' or we need to map it.
            // For now, let's just send what we have, but we might need to adjust based on backend API.

            await api.patch(`/issues/${issueId}`, payload);
        } catch (error) {
            console.error('Failed to update issue', error);
            // Revert on failure (could be improved with a previous state snapshot)
            // For now, just re-fetch issues to ensure sync
            const currentProject = get().currentProject;
            if (currentProject) {
                get().fetchIssues(currentProject.id);
            }
            throw error;
        }
    },

    // Organization & Member Management Actions
    fetchUserRole: async (orgId: string) => {
        try {
            const userId = get().user?.id;
            if (!userId) return;

            const response = await api.get(`/organizations/${orgId}/members/${userId}/role`);
            set({ userRole: response.data.role });
        } catch (error) {
            console.error('Failed to fetch user role', error);
            set({ userRole: null });
        }
    },

    fetchOrganizations: async () => {
        try {
            const userId = get().user?.id;
            if (!userId) return;

            const response = await api.get(`/organizations?user_id=${userId}`);
            set({ organizations: response.data.organizations || [] });
        } catch (error) {
            console.error('Failed to fetch organizations', error);
        }
    },

    fetchOrgMembers: async (orgId: string) => {
        try {
            const response = await api.get(`/organizations/${orgId}/members`);
            set({ orgMembers: response.data.members || [] });
        } catch (error) {
            console.error('Failed to fetch org members', error);
        }
    },

    inviteMember: async (orgId: string, email: string, role: 'admin' | 'member') => {
        try {
            await api.post(`/organizations/${orgId}/invites`, { email, role });
            // Refresh invites list
            await get().fetchInvites(orgId);
        } catch (error) {
            console.error('Failed to invite member', error);
            throw error;
        }
    },

    removeMember: async (orgId: string, userId: string) => {
        try {
            await api.delete(`/organizations/${orgId}/members/${userId}`);
            // Refresh members list
            await get().fetchOrgMembers(orgId);
        } catch (error) {
            console.error('Failed to remove member', error);
            throw error;
        }
    },

    fetchInvites: async (orgId: string) => {
        try {
            const response = await api.get(`/organizations/${orgId}/invites`);
            set({ invites: response.data.invites || [] });
        } catch (error) {
            console.error('Failed to fetch invites', error);
        }
    },

    revokeInvite: async (inviteId: string) => {
        try {
            await api.delete(`/invites/${inviteId}`);
            // Remove from local state
            set((state) => ({
                invites: state.invites.filter(inv => inv.id !== inviteId)
            }));
        } catch (error) {
            console.error('Failed to revoke invite', error);
            throw error;
        }
    },

    addProjectMember: async (projectId: string, userId: string, role: string) => {
        try {
            await api.post(`/projects/${projectId}/members`, { user_id: userId, role });
        } catch (error) {
            console.error('Failed to add project member', error);
            throw error;
        }
    },

    removeProjectMember: async (projectId: string, userId: string) => {
        try {
            await api.delete(`/projects/${projectId}/members/${userId}`);
        } catch (error) {
            console.error('Failed to remove project member', error);
            throw error;
        }
    },
}));
