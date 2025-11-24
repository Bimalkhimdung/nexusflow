import { useState, useEffect } from 'react';
import { useAppStore } from '@/stores';
import { Issue } from '@/types';
import { IssueDetail } from '@/components/issues/IssueDetail';
import { Plus, Search, Filter } from 'lucide-react';

export function Issues() {
    const { issues, fetchIssues, projects } = useAppStore();
    const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);
    const [searchQuery, setSearchQuery] = useState('');

    useEffect(() => {
        if (projects.length > 0) {
            fetchIssues(projects[0].id);
        }
    }, [projects]);

    const filteredIssues = issues.filter(issue =>
        issue.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        issue.key.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="h-full flex flex-col space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-3xl font-bold tracking-tight">Issues</h2>
                    <p className="text-muted-foreground">Track and manage your project tasks.</p>
                </div>
                <button className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2">
                    <Plus className="mr-2 h-4 w-4" />
                    New Issue
                </button>
            </div>

            <div className="flex items-center gap-4">
                <div className="relative flex-1 max-w-sm">
                    <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                    <input
                        type="text"
                        placeholder="Search issues..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 pl-9 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    />
                </div>
                <button className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2">
                    <Filter className="mr-2 h-4 w-4" />
                    Filter
                </button>
            </div>

            <div className="rounded-md border">
                <div className="w-full overflow-auto">
                    <table className="w-full caption-bottom text-sm">
                        <thead className="[&_tr]:border-b">
                            <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground w-[100px]">Key</th>
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Title</th>
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground w-[100px]">Status</th>
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground w-[100px]">Priority</th>
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground w-[150px]">Assignee</th>
                            </tr>
                        </thead>
                        <tbody className="[&_tr:last-child]:border-0">
                            {filteredIssues.map((issue) => (
                                <tr
                                    key={issue.id}
                                    onClick={() => setSelectedIssue(issue)}
                                    className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted cursor-pointer"
                                >
                                    <td className="p-4 align-middle font-mono text-xs">{issue.key}</td>
                                    <td className="p-4 align-middle font-medium">{issue.title}</td>
                                    <td className="p-4 align-middle">
                                        <div className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 ${issue.status === 'DONE' ? 'border-transparent bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' :
                                                issue.status === 'IN_PROGRESS' ? 'border-transparent bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400' :
                                                    'border-transparent bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-400'
                                            }`}>
                                            {issue.status.replace('_', ' ')}
                                        </div>
                                    </td>
                                    <td className="p-4 align-middle">
                                        <div className={`flex items-center gap-1.5 ${issue.priority === 'HIGH' ? 'text-red-500' :
                                                issue.priority === 'MEDIUM' ? 'text-yellow-500' : 'text-blue-500'
                                            }`}>
                                            {issue.priority}
                                        </div>
                                    </td>
                                    <td className="p-4 align-middle">
                                        {issue.assigneeId ? (
                                            <div className="flex items-center gap-2">
                                                <div className="h-6 w-6 rounded-full bg-primary/10 flex items-center justify-center text-[10px] font-bold text-primary">
                                                    {issue.assigneeId.charAt(0)}
                                                </div>
                                            </div>
                                        ) : (
                                            <span className="text-muted-foreground text-xs">Unassigned</span>
                                        )}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {selectedIssue && (
                <IssueDetail issue={selectedIssue} onClose={() => setSelectedIssue(null)} />
            )}
        </div>
    );
}
