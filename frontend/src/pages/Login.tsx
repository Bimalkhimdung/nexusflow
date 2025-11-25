import { useState } from 'react';
import OAuthButton from '../components/auth/OAuthButton';
import EmailPasswordForm from '../components/auth/EmailPasswordForm';

export default function Login() {
    const [isLoading, setIsLoading] = useState(false);

    const handleOAuthLogin = (provider: 'google' | 'github') => {
        setIsLoading(true);
        console.log(`Logging in with ${provider}`);

        setTimeout(() => {
            alert(`OAuth with ${provider} will be implemented in Phase 2`);
            setIsLoading(false);
        }, 1000);
    };

    const handleEmailPasswordSubmit = (email: string, password: string, isSignUp: boolean) => {
        setIsLoading(true);
        console.log(`${isSignUp ? 'Sign up' : 'Sign in'} with:`, { email, password });

        setTimeout(() => {
            if (isSignUp) {
                alert('Sign up successful! Check your email for verification.');
            } else {
                alert('Login successful! Redirecting to dashboard...');
            }
            setIsLoading(false);
        }, 1000);
    };

    return (
        <div className="min-h-screen flex relative overflow-hidden">
            {/* Animated background with darker gradient - no pink */}
            <div className="absolute inset-0 bg-gradient-to-br from-slate-500 via-blue-900 to-black-100 animate-gradient-shift" />

            {/* Floating animated shapes */}
            <div className="absolute inset-0 overflow-hidden pointer-events-none">
                {/* Large floating circles */}
                <div className="absolute top-20 left-10 w-72 h-72 bg-blue-500/10 rounded-full blur-3xl animate-float-slow" />
                <div className="absolute top-40 right-20 w-96 h-96 bg-indigo-400/15 rounded-full blur-3xl animate-float-slower" />
                <div className="absolute bottom-20 left-1/4 w-64 h-64 bg-blue-400/10 rounded-full blur-3xl animate-float" />

                {/* Small floating particles */}
                <div className="absolute top-1/4 left-1/3 w-4 h-4 bg-white/40 rounded-full animate-float-particle" />
                <div className="absolute top-1/3 right-1/4 w-3 h-3 bg-white/30 rounded-full animate-float-particle-delayed" />
                <div className="absolute bottom-1/3 left-1/2 w-5 h-5 bg-white/20 rounded-full animate-float-particle-slow" />
                <div className="absolute top-2/3 right-1/3 w-4 h-4 bg-white/25 rounded-full animate-float-particle" />

                {/* Geometric shapes */}
                <div className="absolute top-1/2 left-20 w-32 h-32 border-2 border-white/20 rotate-45 animate-spin-slow" />
                <div className="absolute bottom-1/4 right-40 w-24 h-24 border-2 border-white/15 rounded-lg animate-spin-reverse" />

                {/* Project Management Objects */}

                {/* Floating Kanban Card 1 */}
                <div className="absolute top-32 right-1/4 animate-float-card">
                    <div className="bg-white/15 backdrop-blur-sm rounded-lg p-3 w-40 border border-white/20 shadow-lg">
                        <div className="h-2 w-16 bg-blue-400/60 rounded mb-2" />
                        <div className="h-1.5 w-full bg-white/30 rounded mb-1" />
                        <div className="h-1.5 w-3/4 bg-white/20 rounded" />
                    </div>
                </div>

                {/* Floating Kanban Card 2 */}
                <div className="absolute bottom-40 left-1/3 animate-float-card-delayed">
                    <div className="bg-white/15 backdrop-blur-sm rounded-lg p-3 w-36 border border-white/20 shadow-lg">
                        <div className="h-2 w-12 bg-green-400/60 rounded mb-2" />
                        <div className="h-1.5 w-full bg-white/30 rounded mb-1" />
                        <div className="h-1.5 w-2/3 bg-white/20 rounded" />
                    </div>
                </div>

                {/* Floating Task Checkmark */}
                <div className="absolute top-1/3 left-1/4 animate-bounce-slow">
                    <div className="w-12 h-12 rounded-full bg-green-500/20 backdrop-blur-sm border border-green-400/30 flex items-center justify-center">
                        <svg className="w-6 h-6 text-green-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                    </div>
                </div>

                {/* Floating Sprint Icon */}
                <div className="absolute bottom-1/3 right-1/4 animate-pulse-slow">
                    <div className="w-14 h-14 rounded-lg bg-purple-500/20 backdrop-blur-sm border border-purple-400/30 flex items-center justify-center">
                        <svg className="w-7 h-7 text-purple-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                    </div>
                </div>

                {/* Floating User Avatar Group */}
                <div className="absolute top-1/2 right-1/3 animate-float-slow">
                    <div className="flex -space-x-2">
                        <div className="w-8 h-8 rounded-full bg-blue-400/40 backdrop-blur-sm border-2 border-white/30" />
                        <div className="w-8 h-8 rounded-full bg-indigo-400/40 backdrop-blur-sm border-2 border-white/30" />
                        <div className="w-8 h-8 rounded-full bg-cyan-400/40 backdrop-blur-sm border-2 border-white/30" />
                    </div>
                </div>

                {/* Floating Calendar/Date */}
                <div className="absolute top-2/3 left-1/4 animate-float-card">
                    <div className="bg-white/15 backdrop-blur-sm rounded-lg p-2 w-16 border border-white/20">
                        <div className="text-center">
                            <div className="text-xs text-white/60 font-semibold">DEC</div>
                            <div className="text-2xl text-white font-bold">25</div>
                        </div>
                    </div>
                </div>

                {/* Floating Progress Bar */}
                <div className="absolute bottom-1/2 left-1/3 animate-pulse-slow">
                    <div className="bg-white/15 backdrop-blur-sm rounded-full p-2 w-32 border border-white/20">
                        <div className="h-2 bg-white/20 rounded-full overflow-hidden">
                            <div className="h-full w-2/3 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-full animate-progress" />
                        </div>
                    </div>
                </div>

                {/* Floating Priority Flag */}
                <div className="absolute top-1/4 right-1/3 animate-wave">
                    <svg className="w-10 h-10 text-red-400/60" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M14.4 6L14 4H5v17h2v-7h5.6l.4 2h7V6z" />
                    </svg>
                </div>

                {/* Floating Comment Bubble */}
                <div className="absolute bottom-1/4 left-1/4 animate-bounce-slow">
                    <div className="bg-white/15 backdrop-blur-sm rounded-lg p-2 border border-white/20">
                        <svg className="w-6 h-6 text-white/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                        </svg>
                    </div>
                </div>

                {/* Floating Tag/Label */}
                <div className="absolute top-1/2 left-1/2 animate-float-particle">
                    <div className="bg-yellow-400/20 backdrop-blur-sm rounded-full px-3 py-1 border border-yellow-400/30">
                        <span className="text-xs text-yellow-200 font-semibold">URGENT</span>
                    </div>
                </div>
            </div>

            {/* Left Side - Branding */}
            <div className="hidden lg:flex lg:w-1/2 p-12 flex-col justify-between relative z-10">
                {/* Logo with pulse effect */}
                <div className="animate-fade-in-up">
                    <div className="inline-block">
                        <h1 className="text-6xl font-bold text-white mb-4 animate-glow">
                            ProjectMgmt
                        </h1>
                        <div className="h-1 w-32 bg-gradient-to-r from-blue-400 to-transparent rounded-full animate-pulse" />
                    </div>
                    <p className="text-2xl text-blue-100 mt-6 animate-fade-in-up animation-delay-200">
                        Manage Your Project with ease
                    </p>
                </div>

                {/* Animated Features */}
                <div className="space-y-6 animate-fade-in-up animation-delay-400">
                    <div className="flex items-start gap-4 group hover:translate-x-2 transition-transform duration-300">
                        <div className="w-14 h-14 rounded-2xl bg-white/20 backdrop-blur-sm flex items-center justify-center flex-shrink-0 group-hover:bg-white/30 transition-all duration-300 group-hover:scale-110">
                            <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                        </div>
                        <div>
                            <h3 className="text-white font-semibold text-xl mb-1">Self-Hosted</h3>
                            <p className="text-blue-100">Own your data, control your workflow</p>
                        </div>
                    </div>

                    <div className="flex items-start gap-4 group hover:translate-x-2 transition-transform duration-300">
                        <div className="w-14 h-14 rounded-2xl bg-white/20 backdrop-blur-sm flex items-center justify-center flex-shrink-0 group-hover:bg-white/30 transition-all duration-300 group-hover:scale-110">
                            <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                            </svg>
                        </div>
                        <div>
                            <h3 className="text-white font-semibold text-xl mb-1">Lightning Fast</h3>
                            <p className="text-blue-100">Built with modern tech for speed</p>
                        </div>
                    </div>

                    <div className="flex items-start gap-4 group hover:translate-x-2 transition-transform duration-300">
                        <div className="w-14 h-14 rounded-2xl bg-white/20 backdrop-blur-sm flex items-center justify-center flex-shrink-0 group-hover:bg-white/30 transition-all duration-300 group-hover:scale-110">
                            <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                            </svg>
                        </div>
                        <div>
                            <h3 className="text-white font-semibold text-xl mb-1">Secure</h3>
                            <p className="text-blue-100">Enterprise-grade security built-in</p>
                        </div>
                    </div>
                </div>

                {/* Footer with fade effect */}
                <div className="text-white/80 text-sm animate-fade-in animation-delay-600">
                    © 2025 NexusFlow. Open Source Project Management.
                </div>
            </div>

            {/* Right Side - Login Form */}
            <div className="flex-1 flex items-center justify-center p-8 relative z-10">
                <div className="w-full max-w-md">
                    {/* Glass morphism card with scale animation */}
                    <div className="backdrop-blur-2xl bg-white/10 border border-white/20 rounded-3xl shadow-2xl p-8 animate-scale-in hover:shadow-3xl transition-shadow duration-500">
                        {/* Mobile Logo */}
                        <div className="lg:hidden text-center mb-8 animate-fade-in">
                            <h1 className="text-4xl font-bold text-white mb-2 animate-glow">
                                NexusFlow
                            </h1>
                            <p className="text-white/80">
                                The open-source Jira killer you can own
                            </p>
                        </div>

                        {/* Welcome Text with slide animation */}
                        <div className="mb-8 animate-slide-in-right">
                            <h2 className="text-3xl font-bold text-white mb-2">Welcome back</h2>
                            <p className="text-white/70">
                                Sign in to your account to continue
                            </p>
                        </div>

                        {/* OAuth Buttons with stagger animation */}
                        <div className="space-y-3 mb-6">
                            <div className="animate-slide-in-left animation-delay-100">
                                <OAuthButton
                                    provider="google"
                                    onClick={() => handleOAuthLogin('google')}
                                    isLoading={isLoading}
                                />
                            </div>
                            <div className="animate-slide-in-left animation-delay-200">
                                <OAuthButton
                                    provider="github"
                                    onClick={() => handleOAuthLogin('github')}
                                    isLoading={isLoading}
                                />
                            </div>
                        </div>

                        {/* Animated Divider */}
                        <div className="relative my-6 animate-fade-in animation-delay-300">
                            <div className="absolute inset-0 flex items-center">
                                <div className="w-full border-t border-white/20" />
                            </div>
                            <div className="relative flex justify-center text-sm">
                                <span className="px-4 bg-gradient-to-r from-transparent via-blue-600/50 to-transparent text-white/80">
                                    Or continue with email
                                </span>
                            </div>
                        </div>

                        {/* Email/Password Form with fade animation */}
                        <div className="animate-fade-in animation-delay-400">
                            <EmailPasswordForm
                                onSubmit={handleEmailPasswordSubmit}
                                isLoading={isLoading}
                            />
                        </div>

                        {/* Footer Links with fade animation */}
                        <div className="mt-8 text-center text-sm text-white/60 animate-fade-in animation-delay-500">
                            <div className="flex items-center justify-center gap-4">
                                <a
                                    href="https://docs.nexusflow.io"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="hover:text-white transition-colors"
                                >
                                    Documentation
                                </a>
                                <span>•</span>
                                <a
                                    href="https://nexusflow.io/privacy"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="hover:text-white transition-colors"
                                >
                                    Privacy Policy
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
