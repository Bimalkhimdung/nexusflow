import { useEffect } from 'react';
import Confetti from 'react-confetti';
import { useNavigate } from 'react-router-dom';
import { CheckCircle, Rocket } from 'lucide-react';

interface Step5Props {
    organizationName: string;
    projectName: string;
}

export default function Step5Complete({ organizationName, projectName }: Step5Props) {
    const navigate = useNavigate();

    useEffect(() => {
        // Scroll to top when component mounts
        window.scrollTo(0, 0);
    }, []);

    const handleGoToDashboard = () => {
        navigate('/dashboard');
    };

    return (
        <div className="animate-fade-in text-center">
            {/* Confetti Animation */}
            <Confetti
                width={window.innerWidth}
                height={window.innerHeight}
                recycle={false}
                numberOfPieces={500}
                gravity={0.3}
            />

            {/* Success Icon */}
            <div className="mb-8">
                <div className="inline-flex items-center justify-center w-24 h-24 bg-green-100 rounded-full mb-6 animate-bounce-slow">
                    <CheckCircle className="w-16 h-16 text-green-600" />
                </div>
                <h2 className="text-4xl font-bold text-gray-900 mb-4">
                    🎉 Welcome to NexusFlow!
                </h2>
                <p className="text-xl text-gray-600">
                    Your workspace is ready to go
                </p>
            </div>

            {/* Summary */}
            <div className="max-w-md mx-auto mb-8">
                <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl p-6 border border-blue-200">
                    <h3 className="font-semibold text-gray-900 mb-4">What we've set up for you:</h3>
                    <div className="space-y-3 text-left">
                        <div className="flex items-start gap-3">
                            <div className="w-6 h-6 rounded-full bg-green-500 flex items-center justify-center flex-shrink-0 mt-0.5">
                                <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                </svg>
                            </div>
                            <div>
                                <p className="font-medium text-gray-900">Organization</p>
                                <p className="text-sm text-gray-600">{organizationName}</p>
                            </div>
                        </div>

                        <div className="flex items-start gap-3">
                            <div className="w-6 h-6 rounded-full bg-green-500 flex items-center justify-center flex-shrink-0 mt-0.5">
                                <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                </svg>
                            </div>
                            <div>
                                <p className="font-medium text-gray-900">First Project</p>
                                <p className="text-sm text-gray-600">{projectName}</p>
                            </div>
                        </div>

                        <div className="flex items-start gap-3">
                            <div className="w-6 h-6 rounded-full bg-green-500 flex items-center justify-center flex-shrink-0 mt-0.5">
                                <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                </svg>
                            </div>
                            <div>
                                <p className="font-medium text-gray-900">Sample Issues</p>
                                <p className="text-sm text-gray-600">8 example tasks to get you started</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Next Steps */}
            <div className="mb-8">
                <h3 className="font-semibold text-gray-900 mb-4">Next steps:</h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-2xl mx-auto text-left">
                    <div className="p-4 bg-white border border-gray-200 rounded-lg">
                        <div className="text-2xl mb-2">📋</div>
                        <h4 className="font-medium text-gray-900 mb-1">Create Issues</h4>
                        <p className="text-sm text-gray-600">Start tracking your work items</p>
                    </div>
                    <div className="p-4 bg-white border border-gray-200 rounded-lg">
                        <div className="text-2xl mb-2">👥</div>
                        <h4 className="font-medium text-gray-900 mb-1">Invite Team</h4>
                        <p className="text-sm text-gray-600">Collaborate with your colleagues</p>
                    </div>
                    <div className="p-4 bg-white border border-gray-200 rounded-lg">
                        <div className="text-2xl mb-2">⚙️</div>
                        <h4 className="font-medium text-gray-900 mb-1">Customize</h4>
                        <p className="text-sm text-gray-600">Tailor workflows to your needs</p>
                    </div>
                </div>
            </div>

            {/* CTA Button */}
            <button
                onClick={handleGoToDashboard}
                className="inline-flex items-center gap-3 px-8 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg font-semibold text-lg hover:from-blue-700 hover:to-indigo-700 transition-all shadow-lg hover:shadow-xl transform hover:scale-105"
            >
                <Rocket className="w-6 h-6" />
                Go to Dashboard
            </button>
        </div>
    );
}
