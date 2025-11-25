interface ProgressIndicatorProps {
    currentStep: number;
    totalSteps: number;
    onStepClick?: (step: number) => void;
}

const steps = [
    { number: 1, label: 'Organization' },
    { number: 2, label: 'Team' },
    { number: 3, label: 'Project' },
    { number: 4, label: 'Git' },
    { number: 5, label: 'Done' },
];

export default function ProgressIndicator({ currentStep, totalSteps, onStepClick }: ProgressIndicatorProps) {
    return (
        <div className="w-full max-w-3xl mx-auto mb-12">
            <div className="flex items-center justify-between">
                {steps.map((step, index) => (
                    <div key={step.number} className="flex items-center flex-1">
                        {/* Step Circle */}
                        <div className="flex flex-col items-center">
                            <button
                                onClick={() => step.number < currentStep && onStepClick?.(step.number)}
                                disabled={step.number >= currentStep}
                                className={`w-12 h-12 rounded-full flex items-center justify-center font-semibold text-sm transition-all duration-300 ${step.number < currentStep
                                        ? 'bg-green-500 text-white cursor-pointer hover:bg-green-600'
                                        : step.number === currentStep
                                            ? 'bg-blue-600 text-white ring-4 ring-blue-200 scale-110'
                                            : 'bg-gray-200 text-gray-400'
                                    }`}
                            >
                                {step.number < currentStep ? (
                                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                                    </svg>
                                ) : (
                                    step.number
                                )}
                            </button>
                            <span className={`mt-2 text-xs font-medium ${step.number === currentStep ? 'text-blue-600' : 'text-gray-500'
                                }`}>
                                {step.label}
                            </span>
                        </div>

                        {/* Connector Line */}
                        {index < steps.length - 1 && (
                            <div className="flex-1 h-1 mx-2 relative">
                                <div className="absolute inset-0 bg-gray-200 rounded" />
                                <div
                                    className={`absolute inset-0 bg-blue-600 rounded transition-all duration-500 ${step.number < currentStep ? 'w-full' : 'w-0'
                                        }`}
                                />
                            </div>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
}
