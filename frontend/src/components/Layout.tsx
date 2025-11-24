import { Link, Outlet, useLocation } from 'react-router-dom';
import { LayoutDashboard, FolderKanban, ListTodo, KanbanSquare, BarChart3, Settings, Bell, Search } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAppStore } from '@/stores';

export function Layout() {
    const location = useLocation();
    const { user } = useAppStore();

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
        </div>
    );
}
