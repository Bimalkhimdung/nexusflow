import { useState } from 'react';
import ProgressIndicator from '../components/onboarding/ProgressIndicator';
import Step1Organization from '../components/onboarding/Step1Organization';
import Step2Invites from '../components/onboarding/Step2Invites';
import Step3Project from '../components/onboarding/Step3Project';
import Step4Git from '../components/onboarding/Step4Git';
import Step5Complete from '../components/onboarding/Step5Complete';

interface OnboardingData {
    organization: {
        name: string;
        slug: string;
        description: string;
        logoUrl?: string;
    };
    invites: Array<{
        email: string;
        role?: 'admin' | 'member';
    }>;
    project: {
        name: string;
        key: string;
        type: 'kanban' | 'scrum' | 'bug-tracking';
        color: string;
    };
    gitRepo?: {
        provider: 'github' | 'gitlab' | 'bitbucket';
        repoUrl: string;
        repoName: string;
    };
}

export default function Onboarding() {
    const [currentStep, setCurrentStep] = useState(1);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const [data, setData] = useState<OnboardingData>({
        organization: {
            name: '',
            slug: '',
            description: '',
        },
        invites: [],
        project: {
            name: '',
            key: '',
            type: 'kanban',
            color: '#3b82f6',
        },
    });

    const handleNext = () => {
        if (currentStep < 5) {
            setCurrentStep(currentStep + 1);
        }
    };

    const handleBack = () => {
        if (currentStep > 1) {
            setCurrentStep(currentStep - 1);
        }
    };

    const handleSkip = () => {
        if (currentStep === 2 || currentStep === 4) {
            handleNext();
        }
    };

    const handleStepClick = (step: number) => {
        setCurrentStep(step);
    };

    const handleOrganizationChange = (orgData: OnboardingData['organization']) => {
        setData({ ...data, organization: orgData });
    };

    const handleInvitesChange = (invites: OnboardingData['invites']) => {
        setData({ ...data, invites });
    };

    const handleProjectChange = (projectData: OnboardingData['project']) => {
        setData({ ...data, project: projectData });
    };

    const handleGitChange = (gitData: OnboardingData['gitRepo']) => {
        setData({ ...data, gitRepo: gitData });
    };

    const handleComplete = async () => {
        setIsSubmitting(true);

        try {
            // TODO: Call API to complete onboarding
            console.log('Onboarding data:', data);

            // Simulate API call
            await new Promise(resolve => setTimeout(resolve, 2000));

            // Move to completion step
            setCurrentStep(5);
        } catch (error) {
            console.error('Failed to complete onboarding:', error);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50 py-12 px-4">
            <div className="max-w-4xl mx-auto">
                {/* Progress Indicator */}
                <ProgressIndicator
                    currentStep={currentStep}
                    totalSteps={5}
                    onStepClick={handleStepClick}
                />

                {/* Step Content */}
                <div className="bg-white rounded-2xl shadow-xl p-8 min-h-[500px]">
                    {currentStep === 1 && (
                        <Step1Organization
                            data={data.organization}
                            onChange={handleOrganizationChange}
                            onNext={handleNext}
                        />
                    )}

                    {currentStep === 2 && (
                        <Step2Invites
                            data={data.invites}
                            onChange={handleInvitesChange}
                            onNext={handleNext}
                            onBack={handleBack}
                            onSkip={handleSkip}
                        />
                    )}

                    {currentStep === 3 && (
                        <Step3Project
                            data={data.project}
                            onChange={handleProjectChange}
                            onNext={handleNext}
                            onBack={handleBack}
                        />
                    )}

                    {currentStep === 4 && (
                        <Step4Git
                            data={data.gitRepo}
                            onChange={handleGitChange}
                            onNext={handleComplete}
                            onBack={handleBack}
                            onSkip={handleSkip}
                            isSubmitting={isSubmitting}
                        />
                    )}

                    {currentStep === 5 && (
                        <Step5Complete
                            organizationName={data.organization.name}
                            projectName={data.project.name}
                        />
                    )}
                </div>
            </div>
        </div>
    );
}
