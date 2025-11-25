import { useState } from 'react';
import { Mail, X, UserPlus } from 'lucide-react';

interface Step2Props {
    data: Array<{
        email: string;
        role?: 'admin' | 'member';
    }>;
    onChange: (data: Step2Props['data']) => void;
    onNext: () => void;
    onBack: () => void;
    onSkip: () => void;
}

export default function Step2Invites({ data, onChange, onNext, onBack, onSkip }: Step2Props) {
    const [email, setEmail] = useState('');
    const [invites, setInvites] = useState(data);
    const [error, setError] = useState('');

    const validateEmail = (email: string) => {
        return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
    };

    const handleAddInvite = () => {
        if (!email) {
            setError('Please enter an email address');
            return;
        }

        if (!validateEmail(email)) {
            setError('Please enter a valid email address');
            return;
        }

        if (invites.some(inv => inv.email === email)) {
            setError('This email has already been added');
            return;
        }

        const newInvites = [...invites, { email, role: 'member' as const }];
        setInvites(newInvites);
        setEmail('');
        setError('');
    };

    const handleRemoveInvite = (emailToRemove: string) => {
        const newInvites = invites.filter(inv => inv.email !== emailToRemove);
        setInvites(newInvites);
    };

    const handleNext = () => {
        onChange(invites);
        onNext();
    };

    const handleSkip = () => {
        onChange([]);
        onSkip();
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleAddInvite();
        }
    };

    return (
        <div className="animate-fade-in">
            <div className="mb-8">
                <h2 className="text-3xl font-bold text-gray-900 mb-2">Invite your team</h2>
                <p className="text-gray-600">
                    Add team members to collaborate on projects. You can invite more people later.
                </p>
            </div>

            {/* Email Input */}
            <div className="mb-6">
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                    Email Address
                </label>
                <div className="flex gap-3">
                    <div className="flex-1 relative">
                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                        <input
                            id="email"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            onKeyPress={handleKeyPress}
                            className={`w-full pl-10 pr-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${error ? 'border-red-500' : 'border-gray-300'
                                }`}
                            placeholder="colleague@example.com"
                        />
                    </div>
                    <button
                        onClick={handleAddInvite}
                        className="px-6 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors flex items-center gap-2"
                    >
                        <UserPlus className="w-5 h-5" />
                        Add
                    </button>
                </div>
                {error && (
                    <p className="text-sm text-red-500 mt-2">{error}</p>
                )}
            </div>

            {/* Invites List */}
            {invites.length > 0 && (
                <div className="mb-8">
                    <h3 className="text-sm font-medium text-gray-700 mb-3">
                        Team Members ({invites.length})
                    </h3>
                    <div className="space-y-2">
                        {invites.map((invite) => (
                            <div
                                key={invite.email}
                                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg border border-gray-200"
                            >
                                <div className="flex items-center gap-3">
                                    <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
                                        <Mail className="w-5 h-5 text-blue-600" />
                                    </div>
                                    <div>
                                        <p className="font-medium text-gray-900">{invite.email}</p>
                                        <p className="text-sm text-gray-500 capitalize">{invite.role || 'Member'}</p>
                                    </div>
                                </div>
                                <button
                                    onClick={() => handleRemoveInvite(invite.email)}
                                    className="p-2 text-gray-400 hover:text-red-500 transition-colors"
                                >
                                    <X className="w-5 h-5" />
                                </button>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* Empty State */}
            {invites.length === 0 && (
                <div className="mb-8 p-8 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300 text-center">
                    <UserPlus className="w-12 h-12 text-gray-400 mx-auto mb-3" />
                    <p className="text-gray-600">No team members added yet</p>
                    <p className="text-sm text-gray-500 mt-1">Add email addresses above to invite your team</p>
                </div>
            )}

            {/* Navigation */}
            <div className="flex justify-between">
                <button
                    onClick={onBack}
                    className="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                    Back
                </button>
                <div className="flex gap-3">
                    <button
                        onClick={handleSkip}
                        className="px-6 py-3 text-gray-600 hover:text-gray-900 font-medium transition-colors"
                    >
                        Skip for now
                    </button>
                    <button
                        onClick={handleNext}
                        className="px-8 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
                    >
                        Continue
                    </button>
                </div>
            </div>
        </div>
    );
}
