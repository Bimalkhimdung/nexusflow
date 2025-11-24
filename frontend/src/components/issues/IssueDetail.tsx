import { useState } from 'react';
import { Issue } from '@/types';
import { X, Calendar, User, Tag, AlertCircle, Clock } from 'lucide-react';
import { useAppStore } from '@/stores';

interface IssueDetailProps {
    issue: Issue;
    onClose: () => void;
}

export function IssueDetail({ issue, onClose }: IssueDetailProps) {
    const { user } = useAppStore();
    const [comment, setComment] = useState('');

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-end bg-black/50 backdrop-blur-sm" onClick={onClose}>
            <div className="h-full w-full max-w-2xl bg-background shadow-2xl animate-in slide-in-from-right duration-300" onClick={e => e.stopPropagation()}>
                <div className="flex h-full flex-col">
                    {/* Header */}
                    <div className="flex items-center justify-between border-b px-6 py-4">
                        <div className="flex items-center gap-3">
                            <span className="font-mono text-sm text-muted-foreground">{issue.key}</span>
                            <div className={`rounded-full px-2.5 py-0.5 text-xs font-semibold ${issue.status === 'DONE' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' :
                                issue.status === 'IN_PROGRESS' ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400' :
                                    'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-400'
                                }`}>
                                {issue.status.replace('_', ' ')}
                            </div>
                        </div>
                        <div className="flex items-center gap-2">
                            <button onClick={onClose} className="rounded-full p-2 hover:bg-secondary">
                                <X className="h-5 w-5" />
                            </button>
                        </div>
                    </div>

                    {/* Content */}
                    <div className="flex-1 overflow-y-auto p-6">
                        <h1 className="text-2xl font-bold mb-6">{issue.title}</h1>

                        <div className="grid grid-cols-3 gap-8">
                            <div className="col-span-2 space-y-8">
                                <div>
                                    <h3 className="text-sm font-medium text-muted-foreground mb-2">Description</h3>
                                    <div className="prose dark:prose-invert max-w-none">
                                        <p>{issue.description}</p>
                                    </div>
                                </div>

                                <div>
                                    <h3 className="text-sm font-medium text-muted-foreground mb-4">Activity</h3>
                                    <div className="flex gap-4">
                                        <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-xs font-bold text-primary">
                                            {user?.fullName.charAt(0)}
                                        </div>
                                        <div className="flex-1 space-y-4">
                                            <textarea
                                                value={comment}
                                                onChange={(e) => setComment(e.target.value)}
                                                placeholder="Add a comment..."
                                                className="w-full rounded-md border bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 min-h-[100px]"
                                            />
                                            <div className="flex justify-end">
                                                <button
                                                    disabled={!comment.trim()}
                                                    className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 py-2"
                                                >
                                                    Comment
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="space-y-6">
                                <div>
                                    <h3 className="text-xs font-medium text-muted-foreground uppercase mb-3">Details</h3>
                                    <div className="space-y-4">
                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                                <User className="h-4 w-4" />
                                                Assignee
                                            </div>
                                            <div className="flex items-center gap-2">
                                                {issue.assigneeId ? (
                                                    <div className="h-6 w-6 rounded-full bg-primary/10 flex items-center justify-center text-[10px] font-bold text-primary">
                                                        {issue.assigneeId.charAt(0)}
                                                    </div>
                                                ) : (
                                                    <span className="text-sm text-muted-foreground">Unassigned</span>
                                                )}
                                            </div>
                                        </div>

                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                                <AlertCircle className="h-4 w-4" />
                                                Priority
                                            </div>
                                            <div className={`flex items-center gap-1.5 text-sm font-medium ${issue.priority === 'HIGH' ? 'text-red-500' :
                                                issue.priority === 'MEDIUM' ? 'text-yellow-500' : 'text-blue-500'
                                                }`}>
                                                {issue.priority}
                                            </div>
                                        </div>

                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                                <Tag className="h-4 w-4" />
                                                Type
                                            </div>
                                            <span className="text-sm">{issue.type}</span>
                                        </div>
                                    </div>
                                </div>

                                <div className="border-t pt-6">
                                    <h3 className="text-xs font-medium text-muted-foreground uppercase mb-3">Dates</h3>
                                    <div className="space-y-4">
                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                                <Calendar className="h-4 w-4" />
                                                Created
                                            </div>
                                            <span className="text-sm">{new Date(issue.createdAt).toLocaleDateString()}</span>
                                        </div>
                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                                <Clock className="h-4 w-4" />
                                                Updated
                                            </div>
                                            <span className="text-sm">{new Date(issue.updatedAt).toLocaleDateString()}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
