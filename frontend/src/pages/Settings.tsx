import { useState } from 'react';
import { User, Bell, Shield, Palette, Globe, Save } from 'lucide-react';

export default function Settings() {
    const [activeTab, setActiveTab] = useState('profile');

    const tabs = [
        { id: 'profile', label: 'Profile', icon: User },
        { id: 'notifications', label: 'Notifications', icon: Bell },
        { id: 'security', label: 'Security', icon: Shield },
        { id: 'appearance', label: 'Appearance', icon: Palette },
        { id: 'general', label: 'General', icon: Globe },
    ];

    return (
        <div className="p-6">
            <div className="mb-6">
                <h1 className="text-3xl font-bold text-gray-900">Settings</h1>
                <p className="text-gray-600 mt-2">Manage your account settings and preferences</p>
            </div>

            <div className="flex gap-6">
                {/* Sidebar */}
                <div className="w-64 flex-shrink-0">
                    <nav className="space-y-1">
                        {tabs.map((tab) => {
                            const Icon = tab.icon;
                            return (
                                <button
                                    key={tab.id}
                                    onClick={() => setActiveTab(tab.id)}
                                    className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left transition-colors ${activeTab === tab.id
                                        ? 'bg-blue-50 text-blue-700 font-medium'
                                        : 'text-gray-700 hover:bg-gray-50'
                                        }`}
                                >
                                    <Icon className="w-5 h-5" />
                                    {tab.label}
                                </button>
                            );
                        })}
                    </nav>
                </div>

                {/* Content */}
                <div className="flex-1 bg-white rounded-lg border border-gray-200 p-6">
                    {activeTab === 'profile' && <ProfileSettings />}
                    {activeTab === 'notifications' && <NotificationSettings />}
                    {activeTab === 'security' && <SecuritySettings />}
                    {activeTab === 'appearance' && <AppearanceSettings />}
                    {activeTab === 'general' && <GeneralSettings />}
                </div>
            </div>
        </div>
    );
}

function ProfileSettings() {
    const [isEditing, setIsEditing] = useState(false);
    const [showConfirmation, setShowConfirmation] = useState(false);
    const [showSuccess, setShowSuccess] = useState(false);
    const [profilePhoto, setProfilePhoto] = useState<string | null>(null);
    const [formData, setFormData] = useState({
        name: 'John Doe',
        email: 'john@example.com',
        bio: 'Product manager passionate about building great software.'
    });

    const handlePhotoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setProfilePhoto(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSaveChanges = () => {
        setShowConfirmation(true);
    };

    const confirmSave = () => {
        // TODO: Save to backend
        console.log('Saving profile:', formData);
        setIsEditing(false);
        setShowConfirmation(false);
        setShowSuccess(true);
    };

    const closeSuccess = () => {
        setShowSuccess(false);
    };

    const cancelSave = () => {
        setShowConfirmation(false);
    };

    return (
        <div>
            <div className="flex items-center justify-between mb-6">
                <h2 className="text-2xl font-bold text-gray-900">Profile Settings</h2>
                {!isEditing && (
                    <button
                        onClick={() => setIsEditing(true)}
                        className="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
                    >
                        Edit Profile
                    </button>
                )}
            </div>

            <div className="space-y-6">
                {/* Avatar - Centered */}
                <div className="flex flex-col items-center">
                    <label className="block text-sm font-medium text-gray-700 mb-4">
                        Profile Picture
                    </label>
                    <div className="w-32 h-32 rounded-full bg-blue-100 flex items-center justify-center overflow-hidden mb-4">
                        {profilePhoto ? (
                            <img src={profilePhoto} alt="Profile" className="w-full h-full object-cover" />
                        ) : (
                            <User className="w-16 h-16 text-blue-600" />
                        )}
                    </div>
                    <div className="text-center">
                        <input
                            id="photo-upload"
                            type="file"
                            accept="image/*"
                            onChange={handlePhotoChange}
                            className="hidden"
                            disabled={!isEditing}
                        />
                        <label
                            htmlFor="photo-upload"
                            className={`inline-block px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium transition-colors ${isEditing
                                    ? 'hover:bg-gray-50 cursor-pointer'
                                    : 'opacity-50 cursor-not-allowed'
                                }`}
                        >
                            Update Image
                        </label>
                        <p className="text-xs text-gray-500 mt-2">JPG, PNG or GIF. Max 2MB.</p>
                    </div>
                </div>

                {/* Name */}
                <div>
                    <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
                        Full Name
                    </label>
                    <input
                        id="name"
                        type="text"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        disabled={!isEditing}
                        className={`w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${!isEditing ? 'bg-gray-50 cursor-not-allowed' : ''
                            }`}
                    />
                </div>

                {/* Email */}
                <div>
                    <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                        Email Address
                    </label>
                    <input
                        id="email"
                        type="email"
                        value={formData.email}
                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        disabled={!isEditing}
                        className={`w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${!isEditing ? 'bg-gray-50 cursor-not-allowed' : ''
                            }`}
                    />
                </div>

                {/* Bio */}
                <div>
                    <label htmlFor="bio" className="block text-sm font-medium text-gray-700 mb-2">
                        Bio
                    </label>
                    <textarea
                        id="bio"
                        rows={4}
                        value={formData.bio}
                        onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                        disabled={!isEditing}
                        className={`w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${!isEditing ? 'bg-gray-50 cursor-not-allowed' : ''
                            }`}
                        placeholder="Tell us about yourself..."
                    />
                </div>

                {/* Action Buttons */}
                {isEditing && (
                    <div className="flex justify-end gap-3">
                        <button
                            onClick={() => setIsEditing(false)}
                            className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                        >
                            Cancel
                        </button>
                        <button
                            onClick={handleSaveChanges}
                            className="flex items-center gap-2 px-6 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
                        >
                            <Save className="w-4 h-4" />
                            Save Changes
                        </button>
                    </div>
                )}
            </div>

            {/* Confirmation Dialog */}
            {showConfirmation && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4 animate-scale-in">
                        <h3 className="text-xl font-bold text-gray-900 mb-4">Confirm Changes</h3>
                        <p className="text-gray-600 mb-6">
                            Are you sure you want to save these changes to your profile?
                        </p>
                        <div className="flex justify-end gap-3">
                            <button
                                onClick={cancelSave}
                                className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={confirmSave}
                                className="px-6 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
                            >
                                Confirm
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Success Dialog */}
            {showSuccess && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4 animate-scale-in">
                        <div className="text-center">
                            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                                <Save className="w-8 h-8 text-green-600" />
                            </div>
                            <h3 className="text-xl font-bold text-gray-900 mb-2">Success!</h3>
                            <p className="text-gray-600 mb-6">
                                Your profile has been updated successfully.
                            </p>
                            <button
                                onClick={closeSuccess}
                                className="px-6 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
                            >
                                OK
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function NotificationSettings() {
    return (
        <div>
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Notification Settings</h2>

            <div className="space-y-4">
                <div className="flex items-center justify-between py-3 border-b">
                    <div>
                        <h3 className="font-medium text-gray-900">Email Notifications</h3>
                        <p className="text-sm text-gray-600">Receive email updates about your projects</p>
                    </div>
                    <input type="checkbox" defaultChecked className="w-5 h-5 text-blue-600" />
                </div>

                <div className="flex items-center justify-between py-3 border-b">
                    <div>
                        <h3 className="font-medium text-gray-900">Push Notifications</h3>
                        <p className="text-sm text-gray-600">Receive push notifications in your browser</p>
                    </div>
                    <input type="checkbox" className="w-5 h-5 text-blue-600" />
                </div>

                <div className="flex items-center justify-between py-3 border-b">
                    <div>
                        <h3 className="font-medium text-gray-900">Issue Updates</h3>
                        <p className="text-sm text-gray-600">Get notified when issues are updated</p>
                    </div>
                    <input type="checkbox" defaultChecked className="w-5 h-5 text-blue-600" />
                </div>

                <div className="flex items-center justify-between py-3">
                    <div>
                        <h3 className="font-medium text-gray-900">Comments</h3>
                        <p className="text-sm text-gray-600">Get notified about new comments</p>
                    </div>
                    <input type="checkbox" defaultChecked className="w-5 h-5 text-blue-600" />
                </div>
            </div>
        </div>
    );
}

function SecuritySettings() {
    return (
        <div>
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Security Settings</h2>

            <div className="space-y-6">
                <div>
                    <h3 className="font-medium text-gray-900 mb-4">Change Password</h3>
                    <div className="space-y-4">
                        <input
                            type="password"
                            placeholder="Current Password"
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                        <input
                            type="password"
                            placeholder="New Password"
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                        <input
                            type="password"
                            placeholder="Confirm New Password"
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                        <button className="px-6 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700">
                            Update Password
                        </button>
                    </div>
                </div>

                <div className="pt-6 border-t">
                    <h3 className="font-medium text-gray-900 mb-4">Two-Factor Authentication</h3>
                    <p className="text-sm text-gray-600 mb-4">
                        Add an extra layer of security to your account
                    </p>
                    <button className="px-6 py-2 border border-gray-300 rounded-lg font-medium hover:bg-gray-50">
                        Enable 2FA
                    </button>
                </div>
            </div>
        </div>
    );
}

function AppearanceSettings() {
    return (
        <div>
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Appearance Settings</h2>

            <div className="space-y-6">
                <div>
                    <h3 className="font-medium text-gray-900 mb-4">Theme</h3>
                    <div className="grid grid-cols-3 gap-4">
                        <button className="p-4 border-2 border-blue-600 rounded-lg bg-white">
                            <div className="w-full h-20 bg-white border border-gray-200 rounded mb-2" />
                            <p className="text-sm font-medium">Light</p>
                        </button>
                        <button className="p-4 border-2 border-gray-200 rounded-lg hover:border-gray-300">
                            <div className="w-full h-20 bg-gray-900 rounded mb-2" />
                            <p className="text-sm font-medium">Dark</p>
                        </button>
                        <button className="p-4 border-2 border-gray-200 rounded-lg hover:border-gray-300">
                            <div className="w-full h-20 bg-gradient-to-br from-white to-gray-900 rounded mb-2" />
                            <p className="text-sm font-medium">Auto</p>
                        </button>
                    </div>
                </div>

                <div>
                    <h3 className="font-medium text-gray-900 mb-4">Accent Color</h3>
                    <div className="flex gap-3">
                        {['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899'].map((color) => (
                            <button
                                key={color}
                                className="w-10 h-10 rounded-full border-2 border-gray-200 hover:scale-110 transition-transform"
                                style={{ backgroundColor: color }}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}

function GeneralSettings() {
    return (
        <div>
            <h2 className="text-2xl font-bold text-gray-900 mb-6">General Settings</h2>

            <div className="space-y-6">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Language
                    </label>
                    <select className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option>English</option>
                        <option>Spanish</option>
                        <option>French</option>
                        <option>German</option>
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Timezone
                    </label>
                    <select className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option>UTC</option>
                        <option>America/New_York</option>
                        <option>Europe/London</option>
                        <option>Asia/Tokyo</option>
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Date Format
                    </label>
                    <select className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option>MM/DD/YYYY</option>
                        <option>DD/MM/YYYY</option>
                        <option>YYYY-MM-DD</option>
                    </select>
                </div>
            </div>
        </div>
    );
}
