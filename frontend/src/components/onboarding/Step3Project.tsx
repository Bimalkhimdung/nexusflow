import { useState, useEffect } from 'react';
import { HexColorPicker } from 'react-colorful';
import { Folder, ChevronDown } from 'lucide-react';

interface Step3Props {
    data: {
        name: string;
        key: string;
        type: 'kanban' | 'scrum' | 'bug-tracking';
        color: string;
    };
    onChange: (data: Step3Props['data']) => void;
    onNext: () => void;
    onBack: () => void;
}

const projectTypes = [
    {
        value: 'kanban' as const,
        label: 'Kanban Board',
        description: 'Visualize work with columns and cards',
        icon: '📋',
    },
    {
        value: 'scrum' as const,
        label: 'Scrum Board',
        description: 'Sprint-based agile development',
        icon: '🏃',
    },
    {
        value: 'bug-tracking' as const,
        label: 'Bug Tracking',
        description: 'Track and manage bugs efficiently',
        icon: '🐛',
    },
];

export default function Step3Project({ data, onChange, onNext, onBack }: Step3Props) {
    const [name, setName] = useState(data.name);
    const [key, setKey] = useState(data.key);
    const [type, setType] = useState(data.type);
    const [color, setColor] = useState(data.color);
    const [showColorPicker, setShowColorPicker] = useState(false);
    const [errors, setErrors] = useState<Record<string, string>>({});

    // Auto-generate key from name
    useEffect(() => {
        if (name && !key) {
            const generatedKey = name
                .toUpperCase()
                .replace(/[^A-Z0-9]+/g, '')
                .slice(0, 10);
            setKey(generatedKey);
        }
    }, [name, key]);

    const validate = () => {
        const newErrors: Record<string, string> = {};

        if (!name || name.length < 2) {
            newErrors.name = 'Project name must be at least 2 characters';
        }

        if (!key || key.length < 2) {
            newErrors.key = 'Project key must be at least 2 characters';
        } else if (key.length > 10) {
            newErrors.key = 'Project key must be at most 10 characters';
        } else if (!/^[A-Z0-9]+$/.test(key)) {
            newErrors.key = 'Project key can only contain uppercase letters and numbers';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleNext = () => {
        if (validate()) {
            onChange({ name, key, type, color });
            onNext();
        }
    };

    return (
        <div className="animate-fade-in">
            <div className="mb-8">
                <h2 className="text-3xl font-bold text-gray-900 mb-2">Create your first project</h2>
                <p className="text-gray-600">
                    Projects help you organize your work. You can create more projects later.
                </p>
            </div>

            {/* Project Name */}
            <div className="mb-6">
                <label htmlFor="project-name" className="block text-sm font-medium text-gray-700 mb-2">
                    Project Name *
                </label>
                <input
                    id="project-name"
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className={`w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.name ? 'border-red-500' : 'border-gray-300'
                        }`}
                    placeholder="Product Roadmap"
                />
                {errors.name && (
                    <p className="text-sm text-red-500 mt-1">{errors.name}</p>
                )}
            </div>

            {/* Project Key */}
            <div className="mb-6">
                <label htmlFor="project-key" className="block text-sm font-medium text-gray-700 mb-2">
                    Project Key *
                </label>
                <input
                    id="project-key"
                    type="text"
                    value={key}
                    onChange={(e) => setKey(e.target.value.toUpperCase().replace(/[^A-Z0-9]/g, '').slice(0, 10))}
                    className={`w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.key ? 'border-red-500' : 'border-gray-300'
                        }`}
                    placeholder="PROD"
                    maxLength={10}
                />
                {errors.key && (
                    <p className="text-sm text-red-500 mt-1">{errors.key}</p>
                )}
                <p className="text-xs text-gray-500 mt-1">
                    This will be used as a prefix for issue IDs (e.g., {key || 'PROJ'}-123)
                </p>
            </div>

            {/* Project Type */}
            <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-3">
                    Project Type *
                </label>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {projectTypes.map((projectType) => (
                        <button
                            key={projectType.value}
                            onClick={() => setType(projectType.value)}
                            className={`p-4 border-2 rounded-lg text-left transition-all ${type === projectType.value
                                    ? 'border-blue-600 bg-blue-50'
                                    : 'border-gray-200 hover:border-gray-300'
                                }`}
                        >
                            <div className="text-3xl mb-2">{projectType.icon}</div>
                            <h3 className="font-semibold text-gray-900 mb-1">{projectType.label}</h3>
                            <p className="text-sm text-gray-600">{projectType.description}</p>
                        </button>
                    ))}
                </div>
            </div>

            {/* Color Picker */}
            <div className="mb-8">
                <label className="block text-sm font-medium text-gray-700 mb-3">
                    Project Color
                </label>
                <div className="relative">
                    <button
                        onClick={() => setShowColorPicker(!showColorPicker)}
                        className="flex items-center gap-3 px-4 py-3 border border-gray-300 rounded-lg hover:border-gray-400 transition-colors"
                    >
                        <div
                            className="w-8 h-8 rounded border-2 border-gray-200"
                            style={{ backgroundColor: color }}
                        />
                        <span className="font-mono text-sm">{color}</span>
                        <ChevronDown className="w-4 h-4 ml-auto text-gray-400" />
                    </button>

                    {showColorPicker && (
                        <div className="absolute top-full mt-2 z-10 p-4 bg-white rounded-lg shadow-xl border border-gray-200">
                            <HexColorPicker color={color} onChange={setColor} />
                            <button
                                onClick={() => setShowColorPicker(false)}
                                className="mt-3 w-full px-4 py-2 bg-gray-100 rounded text-sm font-medium hover:bg-gray-200 transition-colors"
                            >
                                Done
                            </button>
                        </div>
                    )}
                </div>
            </div>

            {/* Navigation */}
            <div className="flex justify-between">
                <button
                    onClick={onBack}
                    className="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                    Back
                </button>
                <button
                    onClick={handleNext}
                    className="px-8 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
                >
                    Continue
                </button>
            </div>
        </div>
    );
}
