import { useEffect } from 'react';
import { useAppStore } from '@/stores';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { ArrowUpRight, CheckCircle2, Circle, Clock } from 'lucide-react';

export function Dashboard() {
    const { fetchProjects, fetchIssues, projects, issues } = useAppStore();

    useEffect(() => {
        fetchProjects();
    }, []);

    useEffect(() => {
        if (projects.length > 0) {
            fetchIssues(projects[0].id);
        }
    }, [projects]);

    const stats = [
        { name: 'Total Issues', value: issues.length, icon: Circle, change: '+12%', changeType: 'positive' },
        { name: 'Completed', value: issues.filter(i => i.status === 'DONE').length, icon: CheckCircle2, change: '+4%', changeType: 'positive' },
        { name: 'In Progress', value: issues.filter(i => i.status === 'IN_PROGRESS').length, icon: Clock, change: '-2%', changeType: 'negative' },
        { name: 'Velocity', value: '24', icon: ArrowUpRight, change: '+8%', changeType: 'positive' },
    ];

    const chartData = [
        { name: 'Mon', completed: 4, created: 2 },
        { name: 'Tue', completed: 3, created: 5 },
        { name: 'Wed', completed: 7, created: 3 },
        { name: 'Thu', completed: 5, created: 4 },
        { name: 'Fri', completed: 8, created: 6 },
        { name: 'Sat', completed: 2, created: 1 },
        { name: 'Sun', completed: 1, created: 0 },
    ];

    return (
        <div className="space-y-8">
            <div>
                <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
                <p className="text-muted-foreground">Overview of your project performance.</p>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {stats.map((stat) => (
                    <div key={stat.name} className="rounded-xl border bg-card text-card-foreground shadow p-6">
                        <div className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <span className="text-sm font-medium text-muted-foreground">{stat.name}</span>
                            <stat.icon className="h-4 w-4 text-muted-foreground" />
                        </div>
                        <div className="text-2xl font-bold">{stat.value}</div>
                        <p className={`text-xs ${stat.changeType === 'positive' ? 'text-green-500' : 'text-red-500'}`}>
                            {stat.change} from last week
                        </p>
                    </div>
                ))}
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <div className="col-span-4 rounded-xl border bg-card text-card-foreground shadow">
                    <div className="p-6 flex flex-col space-y-3">
                        <h3 className="font-semibold leading-none tracking-tight">Issue Activity</h3>
                        <p className="text-sm text-muted-foreground">Issues created vs completed over the last 7 days.</p>
                    </div>
                    <div className="p-6 pt-0 pl-2">
                        <ResponsiveContainer width="100%" height={350}>
                            <BarChart data={chartData}>
                                <XAxis dataKey="name" stroke="#888888" fontSize={12} tickLine={false} axisLine={false} />
                                <YAxis stroke="#888888" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(value) => `${value}`} />
                                <Tooltip
                                    contentStyle={{ backgroundColor: 'hsl(var(--card))', borderColor: 'hsl(var(--border))', borderRadius: 'var(--radius)' }}
                                    itemStyle={{ color: 'hsl(var(--foreground))' }}
                                />
                                <Bar dataKey="created" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} name="Created" />
                                <Bar dataKey="completed" fill="hsl(var(--secondary))" radius={[4, 4, 0, 0]} name="Completed" />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </div>

                <div className="col-span-3 rounded-xl border bg-card text-card-foreground shadow">
                    <div className="p-6 flex flex-col space-y-3">
                        <h3 className="font-semibold leading-none tracking-tight">Recent Issues</h3>
                        <p className="text-sm text-muted-foreground">Latest issues updated in your project.</p>
                    </div>
                    <div className="p-6 pt-0">
                        <div className="space-y-4">
                            {issues.slice(0, 5).map((issue) => (
                                <div key={issue.id} className="flex items-center">
                                    <div className={`mr-4 h-2 w-2 rounded-full ${issue.priority === 'HIGH' ? 'bg-red-500' :
                                        issue.priority === 'MEDIUM' ? 'bg-yellow-500' : 'bg-blue-500'
                                        }`} />
                                    <div className="space-y-1">
                                        <p className="text-sm font-medium leading-none">{issue.title}</p>
                                        <p className="text-xs text-muted-foreground">{issue.key} • {issue.status}</p>
                                    </div>
                                    <div className="ml-auto font-medium text-xs text-muted-foreground">
                                        {new Date(issue.updatedAt).toLocaleDateString()}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
