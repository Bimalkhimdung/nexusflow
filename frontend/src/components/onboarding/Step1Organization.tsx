import { useState, useEffect } from 'react';
import { Upload, Building2 } from 'lucide-react';

interface Step1Props {
    data: {
        name: string;
        slug: string;
        description: string;
        logoUrl?: string;
    };
    onChange: (data: Step1Props['data']) => void;
    onNext: () => void;
}

export default function Step1Organization({ data, onChange, onNext }: Step1Props) {
    const [name, setName] = useState(data.name);
    const [slug, setSlug] = useState(data.slug);
    const [description, setDescription] = useState(data.description);
    const [logoPreview, setLogoPreview] = useState<string | undefined>(data.logoUrl);
    const [errors, setErrors] = useState<Record<string, string>>({});

    // Auto-generate slug from name
    useEffect(() => {
        if (name && !slug) {
            const generatedSlug = name
                .toLowerCase()
                .replace(/[^a-z0-9]+/g, '-')
                .replace(/^-|-$/g, '');
            setSlug(generatedSlug);
        }
    }, [name, slug]);

    const handleLogoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            if (file.size > 2 * 1024 * 1024) {
                setErrors({ ...errors, logo: 'File size must be less than 2MB' });
                return;
            }

            const reader = new FileReader();
            reader.onloadend = () => {
                setLogoPreview(reader.result as string);
                setErrors({ ...errors, logo: '' });
            };
            reader.readAsDataURL(file);
        }
    };

    const validate = () => {
        const newErrors: Record<string, string> = {};

        if (!name || name.length < 2) {
            newErrors.name = 'Organization name must be at least 2 characters';
        }

        if (!slug || slug.length < 2) {
            newErrors.slug = 'Slug must be at least 2 characters';
        } else if (!/^[a-z0-9-]+$/.test(slug)) {
            newErrors.slug = 'Slug can only contain lowercase letters, numbers, and hyphens';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleNext = () => {
        if (validate()) {
            onChange({ name, slug, description, logoUrl: logoPreview });
            onNext();
        }
    };

    return (
        <div className="animate-fade-in">
            <div className="mb-8">
                <h2 className="text-3xl font-bold text-gray-900 mb-2">Create your organization</h2>
                <p className="text-gray-600">
                    Let's start by setting up your workspace. You can always change these later.
                </p>
            </div>

            {/* Logo Upload */}
            <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-3">
                    Organization Logo (Optional)
                </label>
                <div className="flex items-center gap-6">
                    <div className="relative">
                        {logoPreview ? (
                            <img
                                src={logoPreview}
                                alt="Organization logo"
                                className="w-24 h-24 rounded-xl object-cover border-2 border-gray-200"
                            />
                        ) : (
                            <div className="w-24 h-24 rounded-xl bg-gradient-to-br from-blue-100 to-indigo-100 flex items-center justify-center border-2 border-gray-200">
                                <Building2 className="w-12 h-12 text-blue-600" />
                            </div>
                        )}
                    </div>
                    <div>
                        <label
                            htmlFor="logo-upload"
                            className="inline-flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 cursor-pointer transition-colors"
                        >
                            <Upload className="w-4 h-4" />
                            Upload Logo
                        </label>
                        <input
                            id="logo-upload"
                            type="file"
                            accept="image/*"
                            onChange={handleLogoChange}
                            className="hidden"
                        />
                        <p className="text-xs text-gray-500 mt-2">PNG, JPG or SVG. Max 2MB.</p>
                        {errors.logo && (
                            <p className="text-xs text-red-500 mt-1">{errors.logo}</p>
                        )}
                    </div>
                </div>
            </div>

            {/* Organization Name */}
            <div className="mb-6">
                <label htmlFor="org-name" className="block text-sm font-medium text-gray-700 mb-2">
                    Organization Name *
                </label>
                <input
                    id="org-name"
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className={`w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.name ? 'border-red-500' : 'border-gray-300'
                        }`}
                    placeholder="Acme Inc."
                />
                {errors.name && (
                    <p className="text-sm text-red-500 mt-1">{errors.name}</p>
                )}
            </div>

            {/* Organization Slug */}
            <div className="mb-6">
                <label htmlFor="org-slug" className="block text-sm font-medium text-gray-700 mb-2">
                    Organization Slug *
                </label>
                <div className="flex items-center gap-2">
                    <span className="text-gray-500">nexusflow.io/</span>
                    <input
                        id="org-slug"
                        type="text"
                        value={slug}
                        onChange={(e) => setSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
                        className={`flex-1 px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors.slug ? 'border-red-500' : 'border-gray-300'
                            }`}
                        placeholder="acme-inc"
                    />
                </div>
                {errors.slug && (
                    <p className="text-sm text-red-500 mt-1">{errors.slug}</p>
                )}
                <p className="text-xs text-gray-500 mt-1">
                    This will be your organization's unique URL
                </p>
            </div>

            {/* Description */}
            <div className="mb-8">
                <label htmlFor="org-description" className="block text-sm font-medium text-gray-700 mb-2">
                    Description (Optional)
                </label>
                <textarea
                    id="org-description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    rows={4}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Tell us about your organization..."
                />
            </div>

            {/* Navigation */}
            <div className="flex justify-end">
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
