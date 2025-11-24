import { useState, useEffect } from 'react';
import { useAppStore } from '../stores';
import { UserPlus, Trash2, Mail } from 'lucide-react';
import InviteMemberModal from '../components/InviteMemberModal';

const ORG_ID = 'e5616b48-3cdf-4812-a54f-fd861d1ff062'; // Hardcoded for demo

export default function Settings() {
    const [activeTab, setActiveTab] = useState('profile');
    const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);

    const { userRole, orgMembers, invites, fetchOrgMembers, fetchInvites, removeMember, revokeInvite } = useAppStore();

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
                                onClick={() => setActiveTab(tab.id)}
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
                        <div>
                            <h2 className="text-xl font-semibold mb-4">Profile Settings</h2>
                            <p className="text-muted-foreground">Manage your profile information.</p>
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
                        <div>
                            <h2 className="text-xl font-semibold mb-4">Organization Settings</h2>
                            <p className="text-muted-foreground">Manage organization details and settings.</p>
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
        </div>
    );
}
