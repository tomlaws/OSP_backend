import { Survey } from '@/types';
import SurveyForm from './form';

async function getSurvey(token: string): Promise<Survey | null> {
  try {
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/surveys/${token}`, {
      cache: 'no-store', // Always fetch fresh
    });
    if (!res.ok) {
        console.error("Fetch failed", res.status, await res.text());
        return null;
    }
    const json = await res.json();
    return json.data;
  } catch (error) {
    console.error('Failed to fetch survey:', error);
    return null;
  }
}

export default async function SurveyPage({
  params,
}: {
  params: Promise<{ token: string }>;
}) {
  const { token } = await params;
  const survey = await getSurvey(token);

  if (!survey) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">404</h1>
          <p className="text-gray-600">Survey not found or invalid token.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white shadow sm:rounded-lg overflow-hidden">
            <div className="px-4 py-5 sm:px-6 bg-indigo-600">
                <h1 className="text-2xl font-bold text-white">{survey.name}</h1>
                <p className="mt-1 text-indigo-100 text-sm">Please answer the following questions.</p>
            </div>
            <div className="px-4 py-5 sm:p-6">
                <SurveyForm survey={survey} />
            </div>
        </div>
      </div>
    </div>
  );
}
