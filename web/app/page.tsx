import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center">
      <div className="bg-white p-8 rounded shadow-lg text-center max-w-md w-full">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Open Survey Platform</h1>
        <p className="text-gray-600 mb-8">
          Welcome to OSP. Create surveys, collect responses, and get AI-powered insights.
        </p>
        <div className="space-y-4">
          <Link
            href="/admin"
            className="block w-full py-3 px-4 rounded-md shadow bg-indigo-600 text-white font-medium hover:bg-indigo-700"
          >
            Go to Admin Dashboard
          </Link>
        </div>
      </div>
    </div>
  );
}
