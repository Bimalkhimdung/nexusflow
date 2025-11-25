export default function Loading() {
    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-blue-50">
            <div className="text-center">
                {/* Spinning loader */}
                <div className="relative w-20 h-20 mx-auto mb-6">
                    <div className="absolute inset-0 border-4 border-blue-200 rounded-full" />
                    <div className="absolute inset-0 border-4 border-blue-600 border-t-transparent rounded-full animate-spin" />
                </div>

                {/* Loading text */}
                <h2 className="text-xl font-semibold text-gray-900 mb-2">Loading NexusFlow...</h2>
                <p className="text-gray-600">Please wait a moment</p>
            </div>
        </div>
    );
}
