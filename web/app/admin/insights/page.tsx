'use client';

import { Suspense, useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { Insight, Survey } from '@/types';

const CONTEXT_TYPES = [
  { value: 'COURSE_FEEDBACK', label: 'Course Feedback' },
  { value: 'PRODUCT_SATISFACTION', label: 'Product Satisfaction' },
  { value: 'EMPLOYEE_ENGAGEMENT', label: 'Employee Engagement' },
  { value: 'EVENT_FEEDBACK', label: 'Event Feedback' },
];

export default function SurveyInsightsPage() {
    return (
        <Suspense fallback={<div className="text-xl font-bold tracking-tight uppercase">Loading_Data...</div>}>
            <SurveyInsightsContent />
        </Suspense>
    );
}

function SurveyInsightsContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const surveyId = searchParams.get('surveyId');

    const [insights, setInsights] = useState<Insight[]>([]);
    const [survey, setSurvey] = useState<Survey | null>(null);
    const [loading, setLoading] = useState(true);
    const [isGenerating, setIsGenerating] = useState(false);
    const [showGenerateModal, setShowGenerateModal] = useState(false);
    const [contextType, setContextType] = useState(CONTEXT_TYPES[0].value);

    // Fetch data
    useEffect(() => {
        if (!surveyId) return;

        const fetchData = async () => {
            setLoading(true);
            try {
                const [insightsRes, surveyRes] = await Promise.all([
                    fetch(`/next-api/admin/insights?surveyId=${surveyId}`),
                    fetch(`/next-api/admin/surveys/${surveyId}`)
                ]);

                if (insightsRes.ok) {
                   const data = await insightsRes.json();
                   setInsights(data.data || []);
                }
                if (surveyRes.ok) {
                    const data = await surveyRes.json();
                    setSurvey(data.data);
                }
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [surveyId]);

    const handleGenerate = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsGenerating(true);
        try {
            const res = await fetch('/next-api/admin/insights', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    survey_id: surveyId,
                    context_type: contextType,
                }),
            });
            if (res.ok) {
                 // Refresh list
                 const insightsRes = await fetch(`/next-api/admin/insights?surveyId=${surveyId}`);
                 const data = await insightsRes.json();
                 setInsights(data.data || []);
                 setShowGenerateModal(false);
            } else {
                alert('Failed to generate insight');
            }
        } catch (e) {
            console.error(e);
            alert('Error');
        } finally {
            setIsGenerating(false);
        }
    };

    if (loading) return <div className="text-xl font-bold tracking-tight uppercase">Loading_Data...</div>;
    if (!surveyId) return <div>Missing survey ID</div>;

    return (
        <div className="relative">
            <div className="md:flex md:items-center md:justify-between mb-8 border-b-4 border-black pb-4">
                <div className="flex-1 min-w-0">
                    <h2 className="text-3xl font-black uppercase tracking-tighter text-black">
                        Insights / {survey ? survey.name : '...'}
                    </h2>
                    <p className="mt-2 text-sm font-mono text-gray-500 uppercase tracking-widest">
                        Manage and generate AI insights for this survey
                    </p>
                </div>
                <div className="mt-4 flex md:mt-0 md:ml-12 self-end space-x-4">
                     <Link
                        href="/admin"
                        className="inline-flex items-center px-4 py-2 border-2 border-black text-sm font-bold uppercase text-black bg-white hover:bg-gray-100 transition-colors"
                    >
                        Back
                    </Link>
                    <button
                        onClick={() => setShowGenerateModal(true)}
                        className="inline-flex items-center px-4 py-2 border-2 border-black text-sm font-bold uppercase text-white bg-black hover:bg-[#D80000] hover:border-[#D80000] transition-colors"
                    >
                        + Generate Analysis
                    </button>
                </div>
            </div>

            <div className="bg-white border-2 border-black">
                <ul className="divide-y-2 divide-black">
                    {insights.length > 0 ? (
                        insights.map((insight) => (
                            <li key={insight.id} className="group hover:bg-gray-50 transition-colors">
                                <div className="px-6 py-4 flex items-center justify-between">
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-center space-x-3">
                                            <span className={`px-2 py-0.5 text-xs font-bold uppercase border border-black ${
                                                insight.status === 'COMPLETED' ? 'bg-green-100 text-green-800' :
                                                insight.status === 'FAILED' ? 'bg-red-100 text-red-800' :
                                                'bg-yellow-100 text-yellow-800'
                                            }`}>
                                                {insight.status}
                                            </span>
                                            <p className="text-lg font-bold text-black uppercase truncate">
                                                {CONTEXT_TYPES.find(c => c.value === insight.context_type)?.label || insight.context_type}
                                            </p>
                                        </div>
                                        <div className="mt-1 flex items-center space-x-6 text-xs font-mono text-gray-500 uppercase">
                                            <span>Created: {new Date(insight.created_at).toLocaleString()}</span>
                                            {insight.completed_at && (
                                                <span>Completed: {new Date(insight.completed_at).toLocaleString()}</span>
                                            )}
                                        </div>
                                    </div>
                                    <div>
                                        <Link
                                            href={`/admin/insights/${insight.id}`}
                                            className="inline-flex items-center px-3 py-1 border-2 border-black text-xs font-bold uppercase text-black hover:bg-black hover:text-white transition-colors"
                                        >
                                            View Details
                                        </Link>
                                    </div>
                                </div>
                            </li>
                        ))
                    ) : (
                        <li className="px-6 py-12 text-center text-gray-500 font-bold uppercase">
                            No insights generated yet.
                        </li>
                    )}
                </ul>
            </div>

            {/* Modal */}
            {showGenerateModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                    <div className="bg-white border-4 border-black p-8 w-full max-w-md shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
                        <h3 className="text-2xl font-black uppercase text-black mb-6">Generate New Insight</h3>
                        <form onSubmit={handleGenerate}>
                            <div className="mb-6">
                                <label htmlFor="context-type" className="block text-sm font-bold uppercase text-gray-700 mb-2">
                                    Analysis Context
                                </label>
                                <select
                                    id="context-type"
                                    className="block w-full border-2 border-black p-2 text-base focus:outline-none focus:ring-2 focus:ring-[#D80000] rounded-none"
                                    value={contextType}
                                    onChange={(e) => setContextType(e.target.value)}
                                >
                                    {CONTEXT_TYPES.map((type) => (
                                        <option key={type.value} value={type.value}>
                                            {type.label}
                                        </option>
                                    ))}
                                </select>
                                <p className="mt-2 text-xs text-gray-500">
                                    Select the context that best describes this survey to guide the AI analysis.
                                </p>
                            </div>
                            <div className="flex justify-end space-x-3">
                                <button
                                    type="button"
                                    onClick={() => setShowGenerateModal(false)}
                                    className="px-4 py-2 border-2 border-black text-sm font-bold uppercase hover:bg-gray-100"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={isGenerating}
                                    className="px-4 py-2 border-2 border-black bg-black text-white text-sm font-bold uppercase hover:bg-[#D80000] hover:border-[#D80000] disabled:opacity-50"
                                >
                                    {isGenerating ? 'Starting...' : 'Start Generation'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
