import { useState } from 'react';
import { Eye, EyeOff, Mail, Lock } from 'lucide-react';

interface EmailPasswordFormProps {
    onSubmit: (email: string, password: string, isSignUp: boolean) => void;
    isLoading?: boolean;
}

export default function EmailPasswordForm({ onSubmit, isLoading }: EmailPasswordFormProps) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isSignUp, setIsSignUp] = useState(false);
    const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

    const validateForm = () => {
        const newErrors: { email?: string; password?: string } = {};

        // Email validation
        if (!email) {
            newErrors.email = 'Email is required';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
            newErrors.email = 'Invalid email address';
        }

        // Password validation
        if (!password) {
            newErrors.password = 'Password is required';
        } else if (isSignUp && password.length < 8) {
            newErrors.password = 'Password must be at least 8 characters';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (validateForm()) {
            onSubmit(email, password, isSignUp);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {/* Email Input */}
            <div>
                <label htmlFor="email" className="block text-sm font-medium mb-2">
                    Email
                </label>
                <div className="relative">
                    <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
                    <input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className={`w-full pl-10 pr-4 py-3 rounded-lg border bg-muted/30 transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:bg-background text-foreground placeholder:text-muted-foreground ${errors.email ? 'border-red-500' : 'border-border'
                            }`}
                        placeholder="you@example.com"
                        disabled={isLoading}
                    />
                </div>
                {errors.email && (
                    <p className="mt-1 text-sm text-red-500">{errors.email}</p>
                )}
            </div>

            {/* Password Input */}
            <div>
                <label htmlFor="password" className="block text-sm font-medium mb-2">
                    Password
                </label>
                <div className="relative">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
                    <input
                        id="password"
                        type={showPassword ? 'text' : 'password'}
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className={`w-full pl-10 pr-12 py-3 rounded-lg border bg-muted/30 transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:bg-background text-foreground placeholder:text-muted-foreground ${errors.password ? 'border-red-500' : 'border-border'
                            }`}
                        placeholder="••••••••"
                        disabled={isLoading}
                    />
                    <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                        disabled={isLoading}
                    >
                        {showPassword ? (
                            <EyeOff className="h-5 w-5" />
                        ) : (
                            <Eye className="h-5 w-5" />
                        )}
                    </button>
                </div>
                {errors.password && (
                    <p className="mt-1 text-sm text-red-500">{errors.password}</p>
                )}
            </div>

            {/* Submit Button */}
            <button
                type="submit"
                disabled={isLoading}
                className="w-full bg-primary text-primary-foreground py-3 rounded-lg font-medium hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-sm hover:shadow-md"
            >
                {isLoading ? (
                    <div className="flex items-center justify-center gap-2">
                        <div className="h-5 w-5 border-2 border-current border-t-transparent rounded-full animate-spin" />
                        <span>{isSignUp ? 'Creating account...' : 'Signing in...'}</span>
                    </div>
                ) : (
                    <span>{isSignUp ? 'Create Account' : 'Sign In'}</span>
                )}
            </button>

            {/* Toggle Sign Up / Sign In */}
            <div className="text-center">
                <button
                    type="button"
                    onClick={() => setIsSignUp(!isSignUp)}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                    disabled={isLoading}
                >
                    {isSignUp ? (
                        <>
                            Already have an account?{' '}
                            <span className="text-primary font-medium">Sign In</span>
                        </>
                    ) : (
                        <>
                            Don't have an account?{' '}
                            <span className="text-primary font-medium">Sign Up</span>
                        </>
                    )}
                </button>
            </div>
        </form>
    );
}
