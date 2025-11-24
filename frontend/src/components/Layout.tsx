import { useState, useEffect } from 'react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { LayoutDashboard, FolderKanban, ListTodo, KanbanSquare, BarChart3, Settings, Bell, Search, ChevronDown, Plus, Check } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAppStore } from '@/stores';

export function Layout() {
    const location = useLocation();
    const navigate = useNavigate();
    const { user, projects, currentProject, setCurrentProject, fetchProjects, fetchIssues } = useAppStore();
    const [isProjectModalOpen, setIsProjectModalOpen] = useState(false);

    useEffect(() => {
        fetchProjects();
    }, []);

    const handleProjectChange = (projectId: string) => {
        const project = projects.find(p => p.id === projectId);
        if (project) {
            setCurrentProject(project);
            fetchIssues(project.id);
            setIsProjectModalOpen(false);
        }
    };

    const navigation = [
        { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
        { name: 'Projects', href: '/projects', icon: FolderKanban },
        { name: 'Issues', href: '/issues', icon: ListTodo },
        { name: 'Board', href: '/board', icon: KanbanSquare },
        { name: 'Analytics', href: '/analytics', icon: BarChart3 },
    ];

    return (
        <div className="min-h-screen bg-background font-sans antialiased">
            {/* Sidebar */}
            <div className="fixed inset-y-0 z-50 flex w-72 flex-col">
                <div className="flex grow flex-col gap-y-5 overflow-y-auto border-r bg-card px-6 pb-4">
                    <div className="flex h-16 shrink-0 items-center">
                        <div className="flex items-center gap-2 font-bold text-xl text-primary">
                            <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center text-primary-foreground">
                                N
                            </div>
                            NexusFlow
                        </div>
                    </div>

                    <div className="mb-4">
                        <label className="text-xs font-semibold text-muted-foreground mb-2 block">Current Project</label>
                        <button
                            onClick={() => setIsProjectModalOpen(true)}
                            className="flex w-full items-center justify-between rounded-md border bg-background px-3 py-2 text-sm ring-offset-background hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                        >
                            <span className="truncate">{currentProject?.name || 'Select Project'}</span>
                            <ChevronDown className="h-4 w-4 opacity-50" />
                        </button>
                        <div className="mt-2">
                            <Link to="/projects" className="text-xs text-primary hover:underline flex items-center gap-1">
                                <Plus className="h-3 w-3" /> Create Project
                            </Link>
                        </div>
                    </div>
                    <nav className="flex flex-1 flex-col">
                        <ul role="list" className="flex flex-1 flex-col gap-y-7">
                            <li>
                                <ul role="list" className="-mx-2 space-y-1">
                                    {navigation.map((item) => {
                                        const isActive = location.pathname.startsWith(item.href);
                                        return (
                                            <li key={item.name}>
                                                <Link
                                                    to={item.href}
                                                    className={cn(
                                                        isActive
                                                            ? 'bg-secondary text-primary'
                                                            : 'text-muted-foreground hover:bg-secondary/50 hover:text-primary',
                                                        'group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6'
                                                    )}
                                                >
                                                    <item.icon
                                                        className={cn(
                                                            isActive ? 'text-primary' : 'text-muted-foreground group-hover:text-primary',
                                                            'h-6 w-6 shrink-0'
                                                        )}
                                                        aria-hidden="true"
                                                    />
                                                    {item.name}
                                                </Link>
                                            </li>
                                        );
                                    })}
                                </ul>
                            </li>

                            <li className="mt-auto">
                                <Link
                                    to="/settings"
                                    className="group -mx-2 flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6 text-muted-foreground hover:bg-secondary/50 hover:text-primary"
                                >
                                    <Settings
                                        className="h-6 w-6 shrink-0 text-muted-foreground group-hover:text-primary"
                                        aria-hidden="true"
                                    />
                                    Settings
                                </Link>
                            </li>
                        </ul>
                    </nav>
                </div>
            </div>

            {/* Main content */}
            <div className="pl-72">
                <div className="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b bg-background px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8">
                    <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
                        <form className="relative flex flex-1" action="#" method="GET">
                            <label htmlFor="search-field" className="sr-only">
                                Search
                            </label>
                            <div className="relative w-full max-w-md">
                                <Search
                                    className="pointer-events-none absolute inset-y-0 left-0 h-full w-5 text-muted-foreground ml-2"
                                    aria-hidden="true"
                                />
                                <input
                                    id="search-field"
                                    className="block h-full w-full border-0 py-0 pl-10 pr-0 text-foreground placeholder:text-muted-foreground focus:ring-0 sm:text-sm bg-transparent"
                                    placeholder="Search..."
                                    type="search"
                                    name="search"
                                    style={{ outline: 'none', boxShadow: 'none' }}
                                />
                            </div>
                        </form>
                        <div className="flex items-center gap-x-4 lg:gap-x-6">
                            <button type="button" className="-m-2.5 p-2.5 text-muted-foreground hover:text-foreground">
                                <span className="sr-only">View notifications</span>
                                <Bell className="h-6 w-6" aria-hidden="true" />
                            </button>
                            <div className="hidden lg:block lg:h-6 lg:w-px lg:bg-border" aria-hidden="true" />
                            <div className="flex items-center gap-x-4 lg:gap-x-6">
                                <span className="hidden lg:flex lg:items-center">
                                    <span className="ml-4 text-sm font-semibold leading-6 text-foreground" aria-hidden="true">
                                        {user?.fullName}
                                    </span>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>

                <main className="py-10">
                    <div className="px-4 sm:px-6 lg:px-8">
                        <Outlet />
                    </div>
                </main>
            </div>

            {/* Project Selection Modal */}
            {isProjectModalOpen && (
                <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
                    <div
                        className="w-full max-w-lg overflow-hidden rounded-xl bg-background shadow-2xl animate-in zoom-in-95 duration-200 border border-border"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <div className="border-b px-6 py-4">
                            <h2 className="text-lg font-semibold tracking-tight">Switch Project</h2>
                            <p className="text-sm text-muted-foreground">Select a project to switch context.</p>
                        </div>
                        <div className="max-h-[400px] overflow-y-auto p-2">
                            {projects.map((project) => (
                                <button
                                    key={project.id}
                                    onClick={() => handleProjectChange(project.id)}
                                    className={cn(
                                        "flex w-full items-center justify-between rounded-lg px-4 py-3 text-sm transition-colors hover:bg-accent hover:text-accent-foreground",
                                        currentProject?.id === project.id && "bg-accent text-accent-foreground"
                                    )}
                                >
                                    <div className="flex items-center gap-3">
                                        <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 text-primary font-bold">
                                            {project.key.substring(0, 2)}
                                        </div>
                                        <div className="flex flex-col items-start">
                                            <span className="font-medium">{project.name}</span>
                                            <span className="text-xs text-muted-foreground">{project.key}</span>
                                        </div>
                                    </div>
                                    {currentProject?.id === project.id && (
                                        <Check className="h-4 w-4 text-primary" />
                                    )}
                                </button>
                            ))}
                        </div>
                        <div className="border-t bg-muted/50 px-6 py-4 flex justify-between items-center">
                            <Link
                                to="/projects"
                                onClick={() => setIsProjectModalOpen(false)}
                                className="text-sm font-medium text-primary hover:underline flex items-center gap-2"
                            >
                                <Plus className="h-4 w-4" /> Create New Project
                            </Link>
                            <button
                                onClick={() => setIsProjectModalOpen(false)}
                                className="rounded-md bg-background px-4 py-2 text-sm font-medium border hover:bg-accent transition-colors"
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
