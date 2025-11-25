import { useState } from 'react';
import { Github, GitBranch } from 'lucide-react';

interface Step4Props {
    data?: {
        provider: 'github' | 'gitlab' | 'bitbucket';
        repoUrl: string;
        repoName: string;
    };
    onChange: (data: Step4Props['data']) => void;
    onNext: () => void;
    onBack: () => void;
    onSkip: () => void;
    isSubmitting: boolean;
}

const gitProviders = [
    {
        value: 'github' as const,
        label: 'GitHub',
        icon: <Github className="w-8 h-8" />,
        color: 'bg-gray-900',
    },
    {
        value: 'gitlab' as const,
        label: 'GitLab',
        icon: <GitBranch className="w-8 h-8" />,
        color: 'bg-orange-600',
    },
    {
        value: 'bitbucket' as const,
        label: 'Bitbucket',
        icon: <GitBranch className="w-8 h-8" />,
        color: 'bg-blue-700',
    },
];

export default function Step4Git({ data, onChange, onNext, onBack, onSkip, isSubmitting }: Step4Props) {
    const [selectedProvider, setSelectedProvider] = useState<typeof gitProviders[0]['value'] | null>(
        data?.provider || null
    );

    const handleConnect = (provider: typeof gitProviders[0]['value']) => {
        setSelectedProvider(provider);
        // TODO: Implement OAuth flow in Phase 2
        alert(`Git integration with ${provider} will be implemented in Phase 2`);
    };

    const handleFinish = () => {
        if (selectedProvider) {
            onChange({
                provider: selectedProvider,
                repoUrl: '',
                repoName: '',
            });
        }
        onNext();
    };

    return (
        <div className="animate-fade-in">
            <div className="mb-8">
                <h2 className="text-3xl font-bold text-gray-900 mb-2">Connect your repository</h2>
                <p className="text-gray-600">
                    Link your Git repository to track commits and branches. This step is optional.
                </p>
            </div>

            {/* Git Providers */}
            <div className="mb-8">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {gitProviders.map((provider) => (
                        <button
                            key={provider.value}
                            onClick={() => handleConnect(provider.value)}
                            disabled={isSubmitting}
                            className={`p-6 border-2 rounded-lg transition-all hover:shadow-lg ${selectedProvider === provider.value
                                    ? 'border-blue-600 bg-blue-50'
                                    : 'border-gray-200 hover:border-gray-300'
                                } disabled:opacity-50 disabled:cursor-not-allowed`}
                        >
                            <div className={`w-16 h-16 ${provider.color} rounded-lg flex items-center justify-center text-white mx-auto mb-4`}>
                                {provider.icon}
                            </div>
                            <h3 className="font-semibold text-gray-900 text-center">{provider.label}</h3>
                            {selectedProvider === provider.value && (
                                <p className="text-sm text-green-600 text-center mt-2">✓ Selected</p>
                            )}
                        </button>
                    ))}
                </div>
            </div>

            {/* Info Box */}
            <div className="mb-8 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <h4 className="font-medium text-blue-900 mb-2">Why connect a repository?</h4>
                <ul className="text-sm text-blue-800 space-y-1">
                    <li>• Link commits to issues automatically</li>
                    <li>• Track branch deployments</li>
                    <li>• View code changes in context</li>
                    <li>• Automate workflows with webhooks</li>
                </ul>
            </div>

            {/* Navigation */}
            <div className="flex justify-between">
                <button
                    onClick={onBack}
                    disabled={isSubmitting}
                    className="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors disabled:opacity-50"
                >
                    Back
                </button>
                <div className="flex gap-3">
                    <button
                        onClick={onSkip}
                        disabled={isSubmitting}
                        className="px-6 py-3 text-gray-600 hover:text-gray-900 font-medium transition-colors disabled:opacity-50"
                    >
                        Skip for now
                    </button>
                    <button
                        onClick={handleFinish}
                        disabled={isSubmitting}
                        className="px-8 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors disabled:opacity-50 min-w-[120px]"
                    >
                        {isSubmitting ? 'Setting up...' : 'Finish'}
                    </button>
                </div>
            </div>
        </div>
    );
}
