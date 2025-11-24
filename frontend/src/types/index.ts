export interface User {
    id: string;
    email: string;
    fullName: string;
    role: string;
}

export interface Project {
    id: string;
    name: string;
    key: string;
    description: string;
    ownerId: string;
    createdAt: string;
}

export interface Issue {
    id: string;
    key: string;
    title: string;
    description: string;
    status: 'TODO' | 'IN_PROGRESS' | 'DONE';
    priority: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT';
    type: 'STORY' | 'BUG' | 'TASK' | 'EPIC';
    assigneeId?: string;
    reporterId: string;
    projectId: string;
    sprintId?: string;
    createdAt: string;
    updatedAt: string;
}

export interface Sprint {
    id: string;
    name: string;
    goal: string;
    startDate: string;
    endDate: string;
    status: 'PLANNED' | 'ACTIVE' | 'COMPLETED';
    projectId: string;
}

export interface Comment {
    id: string;
    content: string;
    issueId: string;
    userId: string;
    createdAt: string;
}
