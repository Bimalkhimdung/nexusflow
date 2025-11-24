import { useState, useEffect } from 'react';
import { useAppStore } from '../stores';
import { UserPlus, Trash2, Mail, Upload, User as UserIcon } from 'lucide-react';
import InviteMemberModal from '../components/InviteMemberModal';
import ConfirmationModal from '../components/ConfirmationModal';

const ORG_ID = 'e5616b48-3cdf-4812-a54f-fd861d1ff062'; // Hardcoded for demo

export default function Settings() {
    const [activeTab, setActiveTab] = useState('profile');
    const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);

    // Profile form state
    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [address, setAddress] = useState('');
    const [bio, setBio] = useState('');
    const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
    const [isSaving, setIsSaving] = useState(false);
    const [saveSuccess, setSaveSuccess] = useState(false);
    const [isEditingProfile, setIsEditingProfile] = useState(false);

    // Organization form state
    const [orgName, setOrgName] = useState('');
    const [orgSlug, setOrgSlug] = useState('');
    const [orgDescription, setOrgDescription] = useState('');
    const [orgLogoPreview, setOrgLogoPreview] = useState<string | null>(null);
    const [isSavingOrg, setIsSavingOrg] = useState(false);
    const [saveOrgSuccess, setSaveOrgSuccess] = useState(false);
    const [isEditingOrg, setIsEditingOrg] = useState(false);

    // Confirmation modal state
    const [isProfileConfirmOpen, setIsProfileConfirmOpen] = useState(false);
    const [isOrgConfirmOpen, setIsOrgConfirmOpen] = useState(false);
    const [isTabSwitchConfirmOpen, setIsTabSwitchConfirmOpen] = useState(false);
    const [pendingTab, setPendingTab] = useState<string | null>(null);

    // Track unsaved changes
    const [hasProfileChanges, setHasProfileChanges] = useState(false);
    const [hasOrgChanges, setHasOrgChanges] = useState(false);

    const { user, userRole, organizations, currentOrganization, orgMembers, invites, fetchOrgMembers, fetchInvites, removeMember, revokeInvite } = useAppStore();

    // Initialize profile form with user data
    useEffect(() => {
        if (user) {
            setFullName(user.fullName || '');
            setEmail(user.email || '');
            setPhone(user.phone || '');
            setAddress(user.address || '');
            setBio(user.bio || '');
            setAvatarPreview(user.avatarUrl || null);
        }
    }, [user]);

    // Initialize organization form with current organization data
    useEffect(() => {
        if (currentOrganization) {
            setOrgName(currentOrganization.name || '');
            setOrgSlug(currentOrganization.slug || '');
            setOrgDescription(currentOrganization.description || '');
            setOrgLogoPreview(currentOrganization.logoUrl || null);
        }
    }, [currentOrganization]);

    // Track profile changes
    useEffect(() => {
        if (!user) return;
        const hasChanges =
            fullName !== (user.fullName || '') ||
            email !== (user.email || '') ||
            phone !== (user.phone || '') ||
            address !== (user.address || '') ||
            bio !== (user.bio || '') ||
            avatarPreview !== (user.avatarUrl || null);
        setHasProfileChanges(hasChanges);
    }, [fullName, email, phone, address, bio, avatarPreview, user]);

    // Track organization changes
    useEffect(() => {
        if (!currentOrganization) return;
        const hasChanges =
            orgName !== (currentOrganization.name || '') ||
            orgSlug !== (currentOrganization.slug || '') ||
            orgDescription !== (currentOrganization.description || '') ||
            orgLogoPreview !== (currentOrganization.logoUrl || null);
        setHasOrgChanges(hasChanges);
    }, [orgName, orgSlug, orgDescription, orgLogoPreview, currentOrganization]);

    useEffect(() => {
        if (activeTab === 'members') {
            fetchOrgMembers(ORG_ID);
        } else if (activeTab === 'invites') {
            fetchInvites(ORG_ID);
        }
    }, [activeTab, fetchOrgMembers, fetchInvites]);

    const isAdmin = userRole === 'admin' || userRole === 'owner';

    const tabs = [
        { id: 'profile', label: 'Profile', adminOnly: false },
        { id: 'appearance', label: 'Appearance', adminOnly: false },
        { id: 'notifications', label: 'Notifications', adminOnly: false },
        { id: 'organization', label: 'Organization', adminOnly: true },
        { id: 'members', label: 'Members', adminOnly: true },
        { id: 'invites', label: 'Invites', adminOnly: true },
    ];

    const visibleTabs = tabs.filter(tab => !tab.adminOnly || isAdmin);

    const handleRemoveMember = async (userId: string) => {
        if (confirm('Are you sure you want to remove this member?')) {
            try {
                await removeMember(ORG_ID, userId);
            } catch (error) {
                console.error('Failed to remove member', error);
            }
        }
    };

    const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setAvatarPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSaveProfile = async () => {
        // Close confirmation modal
        setIsProfileConfirmOpen(false);

        setIsSaving(true);
        setSaveSuccess(false);

        try {
            // TODO: Call API to update user profile
            // await api.patch('/users/me', { fullName, email, phone, address, bio, avatarUrl: avatarPreview });

            // Simulate API call
            await new Promise(resolve => setTimeout(resolve, 1000));

            setSaveSuccess(true);
            setIsEditingProfile(false); // Exit edit mode after successful save
            setHasProfileChanges(false); // Reset changes flag
            setTimeout(() => setSaveSuccess(false), 3000);
        } catch (error) {
            console.error('Failed to save profile', error);
        } finally {
            setIsSaving(false);
        }
    };

    const handleCancelEdit = () => {
        // Reset form to original values
        if (user) {
            setFullName(user.fullName || '');
            setEmail(user.email || '');
            setPhone(user.phone || '');
            setAddress(user.address || '');
            setBio(user.bio || '');
            setAvatarPreview(user.avatarUrl || null);
        }
        setIsEditingProfile(false);
        setHasProfileChanges(false); // Reset changes flag
    };

    const handleOrgLogoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setOrgLogoPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSaveOrganization = async () => {
        // Close confirmation modal
        setIsOrgConfirmOpen(false);

        setIsSavingOrg(true);
        setSaveOrgSuccess(false);

        try {
            // TODO: Call API to update organization
            // if (currentOrganization) {
            //   await api.patch(`/organizations/${currentOrganization.id}`, {
            //     name: orgName,
            //     slug: orgSlug,
            //     description: orgDescription,
            //     logoUrl: orgLogoPreview
            //   });
            // } else {
            //   await api.post('/organizations', {
            //     name: orgName,
            //     slug: orgSlug,
            //     description: orgDescription,
            //     logoUrl: orgLogoPreview
            //   });
            // }

            // Simulate API call
            await new Promise(resolve => setTimeout(resolve, 1000));

            setSaveOrgSuccess(true);
            setIsEditingOrg(false); // Exit edit mode after successful save
            setHasOrgChanges(false); // Reset changes flag
            setTimeout(() => setSaveOrgSuccess(false), 3000);
        } catch (error) {
            console.error('Failed to save organization', error);
        } finally {
            setIsSavingOrg(false);
        }
    };

    const handleCancelOrgEdit = () => {
        // Reset form to original values
        if (currentOrganization) {
            setOrgName(currentOrganization.name || '');
            setOrgSlug(currentOrganization.slug || '');
            setOrgDescription(currentOrganization.description || '');
            setOrgLogoPreview(currentOrganization.logoUrl || null);
        }
        setIsEditingOrg(false);
        setHasOrgChanges(false); // Reset changes flag
    };

    // Handle tab switching with unsaved changes check
    const handleTabSwitch = (tab: string) => {
        // Check if we're in edit mode with unsaved changes
        const inProfileEdit = activeTab === 'profile' && isEditingProfile && hasProfileChanges;
        const inOrgEdit = activeTab === 'organization' && isEditingOrg && hasOrgChanges;

        if (inProfileEdit || inOrgEdit) {
            // Store the tab they want to switch to
            setPendingTab(tab);
            setIsTabSwitchConfirmOpen(true);
        } else {
            // No unsaved changes, switch immediately
            setActiveTab(tab);
        }
    };

    // Confirm tab switch and discard changes
    const confirmTabSwitch = () => {
        if (pendingTab) {
            // Reset form to original values
            if (activeTab === 'profile') {
                handleCancelEdit();
            } else if (activeTab === 'organization') {
                handleCancelOrgEdit();
            }

            setActiveTab(pendingTab);
            setPendingTab(null);
        }
        setIsTabSwitchConfirmOpen(false);
    };

    // Cancel tab switch
    const cancelTabSwitch = () => {
        setPendingTab(null);
        setIsTabSwitchConfirmOpen(false);
    };

    const handleRevokeInvite = async (inviteId: string) => {
        if (confirm('Are you sure you want to revoke this invite?')) {
            try {
                await revokeInvite(inviteId);
            } catch (error) {
                console.error('Failed to revoke invite', error);
            }
        }
    };

    return (
        <div className="min-h-screen bg-background">
            <div className="max-w-6xl mx-auto p-6">
                <h1 className="text-3xl font-bold mb-6">Settings</h1>

                {/* Tabs */}
                <div className="border-b border-border mb-6">
                    <div className="flex gap-4">
                        {visibleTabs.map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() => handleTabSwitch(tab.id)}
                                className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${activeTab === tab.id
                                    ? 'border-primary text-primary'
                                    : 'border-transparent text-muted-foreground hover:text-foreground'
                                    }`}
                            >
                                {tab.label}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Tab Content */}
                <div className="bg-card rounded-lg border border-border p-6">
                    {activeTab === 'profile' && (
                        <div className="max-w-2xl">
                            <div className="flex items-center justify-between mb-6">
                                <h2 className="text-xl font-semibold">Profile Settings</h2>
                                {!isEditingProfile && (
                                    <button
                                        onClick={() => setIsEditingProfile(true)}
                                        className="inline-flex items-center gap-2 rounded-md border border-input bg-background px-4 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
                                    >
                                        Edit Profile
                                    </button>
                                )}
                            </div>

                            {/* Avatar Upload */}
                            <div className="mb-6">
                                <label className="block text-sm font-medium mb-3">Profile Picture</label>
                                <div className="flex items-center gap-4">
                                    <div className="relative">
                                        {avatarPreview ? (
                                            <img
                                                src={avatarPreview}
                                                alt="Profile"
                                                className="h-20 w-20 rounded-full object-cover border-2 border-border"
                                            />
                                        ) : (
                                            <div className="h-20 w-20 rounded-full bg-secondary flex items-center justify-center border-2 border-border">
                                                <UserIcon className="h-10 w-10 text-muted-foreground" />
                                            </div>
                                        )}
                                    </div>
                                    {isEditingProfile && (
                                        <div>
                                            <label
                                                htmlFor="avatar-upload"
                                                className="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 cursor-pointer transition-colors"
                                            >
                                                <Upload className="h-4 w-4" />
                                                Upload Photo
                                            </label>
                                            <input
                                                id="avatar-upload"
                                                type="file"
                                                accept="image/*"
                                                onChange={handleAvatarChange}
                                                className="hidden"
                                            />
                                            <p className="text-xs text-muted-foreground mt-2">JPG, PNG or GIF. Max 2MB.</p>
                                        </div>
                                    )}
                                </div>
                            </div>

                            {/* Full Name */}
                            <div className="mb-4">
                                <label htmlFor="fullName" className="block text-sm font-medium mb-2">
                                    Full Name
                                </label>
                                <input
                                    id="fullName"
                                    type="text"
                                    value={fullName}
                                    onChange={(e) => setFullName(e.target.value)}
                                    disabled={!isEditingProfile}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="John Doe"
                                />
                            </div>

                            {/* Email */}
                            <div className="mb-4">
                                <label htmlFor="email" className="block text-sm font-medium mb-2">
                                    Email Address
                                </label>
                                <input
                                    id="email"
                                    type="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    disabled={!isEditingProfile}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="john@example.com"
                                />
                            </div>

                            {/* Phone */}
                            <div className="mb-4">
                                <label htmlFor="phone" className="block text-sm font-medium mb-2">
                                    Phone Number
                                </label>
                                <input
                                    id="phone"
                                    type="tel"
                                    value={phone}
                                    onChange={(e) => setPhone(e.target.value)}
                                    disabled={!isEditingProfile}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="+1 (555) 123-4567"
                                />
                            </div>

                            {/* Address */}
                            <div className="mb-4">
                                <label htmlFor="address" className="block text-sm font-medium mb-2">
                                    Address
                                </label>
                                <textarea
                                    id="address"
                                    value={address}
                                    onChange={(e) => setAddress(e.target.value)}
                                    disabled={!isEditingProfile}
                                    rows={3}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="123 Main St, City, State, ZIP"
                                />
                            </div>

                            {/* Bio */}
                            <div className="mb-6">
                                <label htmlFor="bio" className="block text-sm font-medium mb-2">
                                    Bio
                                </label>
                                <textarea
                                    id="bio"
                                    value={bio}
                                    onChange={(e) => setBio(e.target.value)}
                                    disabled={!isEditingProfile}
                                    rows={4}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="Tell us about yourself..."
                                />
                            </div>

                            {/* Success Message */}
                            {saveSuccess && (
                                <div className="mb-4 rounded-md bg-green-500/10 p-3 text-sm text-green-600">
                                    Profile updated successfully!
                                </div>
                            )}

                            {/* Action Buttons */}
                            {isEditingProfile && (
                                <div className="flex justify-end gap-2">
                                    <button
                                        onClick={handleCancelEdit}
                                        className="inline-flex items-center justify-center rounded-md border border-input bg-background px-6 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        onClick={() => setIsProfileConfirmOpen(true)}
                                        disabled={isSaving}
                                        className="inline-flex items-center justify-center rounded-md bg-primary px-6 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:pointer-events-none transition-colors"
                                    >
                                        {isSaving ? 'Saving...' : 'Save Changes'}
                                    </button>
                                </div>
                            )}
                        </div>
                    )}
                    {activeTab === 'appearance' && (
                        <div>
                            <h2 className="text-xl font-semibold mb-4">Appearance</h2>
                            <p className="text-muted-foreground">Customize the look and feel of the application.</p>
                        </div>
                    )}

                    {activeTab === 'notifications' && (
                        <div>
                            <h2 className="text-xl font-semibold mb-4">Notifications</h2>
                            <p className="text-muted-foreground">Manage your notification preferences.</p>
                        </div>
                    )}

                    {activeTab === 'organization' && isAdmin && (
                        <div className="max-w-2xl">
                            <div className="flex items-center justify-between mb-6">
                                <h2 className="text-xl font-semibold">Organization Settings</h2>
                                {!isEditingOrg && (
                                    <button
                                        onClick={() => setIsEditingOrg(true)}
                                        className="inline-flex items-center gap-2 rounded-md border border-input bg-background px-4 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
                                    >
                                        Edit Organization
                                    </button>
                                )}
                            </div>

                            {/* Organization Logo */}
                            <div className="mb-6">
                                <label className="block text-sm font-medium mb-3">Organization Logo</label>
                                <div className="flex items-center gap-4">
                                    <div className="relative">
                                        {orgLogoPreview ? (
                                            <img
                                                src={orgLogoPreview}
                                                alt="Organization Logo"
                                                className="h-20 w-20 rounded-lg object-cover border-2 border-border"
                                            />
                                        ) : (
                                            <div className="h-20 w-20 rounded-lg bg-secondary flex items-center justify-center border-2 border-border">
                                                <UserIcon className="h-10 w-10 text-muted-foreground" />
                                            </div>
                                        )}
                                    </div>
                                    {isEditingOrg && (
                                        <div>
                                            <label
                                                htmlFor="org-logo-upload"
                                                className="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 cursor-pointer transition-colors"
                                            >
                                                <Upload className="h-4 w-4" />
                                                Upload Logo
                                            </label>
                                            <input
                                                id="org-logo-upload"
                                                type="file"
                                                accept="image/*"
                                                onChange={handleOrgLogoChange}
                                                className="hidden"
                                            />
                                            <p className="text-xs text-muted-foreground mt-2">PNG or SVG recommended. Max 2MB.</p>
                                        </div>
                                    )}
                                </div>
                            </div>

                            {/* Organization Name */}
                            <div className="mb-4">
                                <label htmlFor="orgName" className="block text-sm font-medium mb-2">
                                    Organization Name
                                </label>
                                <input
                                    id="orgName"
                                    type="text"
                                    value={orgName}
                                    onChange={(e) => setOrgName(e.target.value)}
                                    disabled={!isEditingOrg}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="Acme Inc."
                                />
                            </div>

                            {/* Organization Slug */}
                            <div className="mb-4">
                                <label htmlFor="orgSlug" className="block text-sm font-medium mb-2">
                                    Organization Slug
                                </label>
                                <div className="flex items-center gap-2">
                                    <span className="text-sm text-muted-foreground">nexusflow.io/</span>
                                    <input
                                        id="orgSlug"
                                        type="text"
                                        value={orgSlug}
                                        onChange={(e) => setOrgSlug(e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, ''))}
                                        disabled={!isEditingOrg}
                                        className="flex-1 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                        placeholder="acme-inc"
                                    />
                                </div>
                                <p className="text-xs text-muted-foreground mt-1">Lowercase letters, numbers, and hyphens only</p>
                            </div>

                            {/* Organization Description */}
                            <div className="mb-6">
                                <label htmlFor="orgDescription" className="block text-sm font-medium mb-2">
                                    Description
                                </label>
                                <textarea
                                    id="orgDescription"
                                    value={orgDescription}
                                    onChange={(e) => setOrgDescription(e.target.value)}
                                    disabled={!isEditingOrg}
                                    rows={4}
                                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                                    placeholder="Tell us about your organization..."
                                />
                            </div>

                            {/* Success Message */}
                            {saveOrgSuccess && (
                                <div className="mb-4 rounded-md bg-green-500/10 p-3 text-sm text-green-600">
                                    Organization updated successfully!
                                </div>
                            )}

                            {/* Action Buttons */}
                            {isEditingOrg && (
                                <div className="flex justify-end gap-2">
                                    <button
                                        onClick={handleCancelOrgEdit}
                                        className="inline-flex items-center justify-center rounded-md border border-input bg-background px-6 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        onClick={() => setIsOrgConfirmOpen(true)}
                                        disabled={isSavingOrg}
                                        className="inline-flex items-center justify-center rounded-md bg-primary px-6 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:pointer-events-none transition-colors"
                                    >
                                        {isSavingOrg ? 'Saving...' : 'Save Changes'}
                                    </button>
                                </div>
                            )}
                        </div>
                    )}
                    {activeTab === 'members' && isAdmin && (
                        <div>
                            <div className="flex items-center justify-between mb-6">
                                <h2 className="text-xl font-semibold">Members</h2>
                                <button
                                    onClick={() => setIsInviteModalOpen(true)}
                                    className="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
                                >
                                    <UserPlus className="h-4 w-4" />
                                    Invite Member
                                </button>
                            </div>

                            <div className="space-y-2">
                                {orgMembers.length === 0 ? (
                                    <p className="text-muted-foreground text-center py-8">No members found</p>
                                ) : (
                                    orgMembers.map((member) => (
                                        <div
                                            key={member.id}
                                            className="flex items-center justify-between p-4 rounded-md border border-border hover:bg-accent/50 transition-colors"
                                        >
                                            <div>
                                                <p className="font-medium">User ID: {member.userId}</p>
                                                <p className="text-sm text-muted-foreground capitalize">{member.role}</p>
                                            </div>
                                            <button
                                                onClick={() => handleRemoveMember(member.userId)}
                                                className="p-2 rounded-md hover:bg-destructive/10 text-destructive transition-colors"
                                            >
                                                <Trash2 className="h-4 w-4" />
                                            </button>
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>
                    )}

                    {activeTab === 'invites' && isAdmin && (
                        <div>
                            <div className="flex items-center justify-between mb-6">
                                <h2 className="text-xl font-semibold">Pending Invites</h2>
                                <button
                                    onClick={() => setIsInviteModalOpen(true)}
                                    className="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
                                >
                                    <Mail className="h-4 w-4" />
                                    New Invite
                                </button>
                            </div>

                            <div className="space-y-2">
                                {invites.length === 0 ? (
                                    <p className="text-muted-foreground text-center py-8">No pending invites</p>
                                ) : (
                                    invites.filter(inv => inv.status === 'pending').map((invite) => (
                                        <div
                                            key={invite.id}
                                            className="flex items-center justify-between p-4 rounded-md border border-border hover:bg-accent/50 transition-colors"
                                        >
                                            <div>
                                                <p className="font-medium">{invite.email}</p>
                                                <p className="text-sm text-muted-foreground capitalize">{invite.role}</p>
                                            </div>
                                            <button
                                                onClick={() => handleRevokeInvite(invite.id)}
                                                className="p-2 rounded-md hover:bg-destructive/10 text-destructive transition-colors"
                                            >
                                                <Trash2 className="h-4 w-4" />
                                            </button>
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>
                    )}
                </div>
            </div>

            <InviteMemberModal
                isOpen={isInviteModalOpen}
                onClose={() => setIsInviteModalOpen(false)}
                orgId={ORG_ID}
            />

            {/* Profile Save Confirmation Modal */}
            <ConfirmationModal
                isOpen={isProfileConfirmOpen}
                onClose={() => setIsProfileConfirmOpen(false)}
                onConfirm={handleSaveProfile}
                title="Save Profile Changes"
                message="Are you sure you want to save these changes to your profile?"
                confirmText="Save Changes"
                isLoading={isSaving}
            />

            {/* Organization Save Confirmation Modal */}
            <ConfirmationModal
                isOpen={isOrgConfirmOpen}
                onClose={() => setIsOrgConfirmOpen(false)}
                onConfirm={handleSaveOrganization}
                title="Save Organization Changes"
                message="Are you sure you want to save these changes to your organization?"
                confirmText="Save Changes"
                isLoading={isSavingOrg}
            />

            {/* Tab Switch Confirmation Modal */}
            <ConfirmationModal
                isOpen={isTabSwitchConfirmOpen}
                onClose={cancelTabSwitch}
                onConfirm={confirmTabSwitch}
                title="Unsaved Changes"
                message="You have unsaved changes. If you leave this page, your changes will be lost. Are you sure you want to continue?"
                confirmText="Discard Changes"
                cancelText="Stay on Page"
            />
        </div>
    );
}
