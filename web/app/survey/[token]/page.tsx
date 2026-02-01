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
      <div className="min-h-screen flex items-center justify-center bg-white">
        <div className="text-center border-2 border-black p-8">
          <h1 className="text-5xl font-black text-black mb-4">404</h1>
          <p className="text-gray-600 font-bold uppercase tracking-widest">Survey not found or invalid token.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      {/* Swiss Header */}
      <div className="bg-black text-white py-4 px-6 md:px-12 mb-12">
        <div className="max-w-4xl mx-auto flex justify-between items-center">
            <h1 className="text-xl font-bold uppercase tracking-widest">
                OSP / Public Survey
            </h1>
            <span className="text-xs font-mono border border-white px-2 py-0.5">SECURE_CONNECTION</span>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-6 md:px-12 pb-24">
        {/* Title Block */}
        <div className="mb-16 border-l-4 border-[#D80000] pl-6">
            <h1 className="text-5xl md:text-6xl font-black text-black tracking-tighter mb-4 leading-none">
                {survey.name}
            </h1>
            <p className="text-gray-500 font-bold uppercase tracking-wide text-sm">
                Please provide your objective responses below.
            </p>
        </div>
        
        <div className="border-t-2 border-black pt-12">
            <SurveyForm survey={survey} />
        </div>
      </div>
    </div>
  );
}
