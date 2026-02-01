'use client';

import { useState, use, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Survey } from '@/types';
import Link from 'next/link';

const CONTEXT_TYPES = [
  { value: 'COURSE_FEEDBACK', label: 'Course Feedback' },
  { value: 'PRODUCT_SATISFACTION', label: 'Product Satisfaction' },
  { value: 'EMPLOYEE_ENGAGEMENT', label: 'Employee Engagement' },
  { value: 'EVENT_FEEDBACK', label: 'Event Feedback' },
];

export default function GenerateInsightPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [contextType, setContextType] = useState(CONTEXT_TYPES[0].value);
  const [loading, setLoading] = useState(false);
  const [fetchingSurvey, setFetchingSurvey] = useState(true);

  useEffect(() => {
    fetch(`/next-api/admin/surveys/${id}`)
      .then(res => res.json())
      .then(data => {
        if (data.data) setSurvey(data.data);
      })
      .catch(console.error)
      .finally(() => setFetchingSurvey(false));
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await fetch('/next-api/admin/insights', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          survey_id: id,
          context_type: contextType,
        }),
      });

      if (res.ok) {
        router.push('/admin');
      } else {
        const errorData = await res.json();
        alert(`Failed to start generation: ${errorData.error || 'Unknown error'}`);
      }
    } catch (e) {
      console.error(e);
      alert('Error generating insight');
    } finally {
      setLoading(false);
    }
  };

  if (fetchingSurvey) return <div>Loading survey details...</div>;
  if (!survey) return <div>Survey not found</div>;

  return (
    <div className="max-w-xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div className="bg-white shadow sm:rounded-lg overflow-hidden">
        <div className="px-4 py-5 sm:px-6 bg-indigo-600">
          <h1 className="text-xl font-bold text-white">Generate Insight</h1>
          <p className="mt-1 text-indigo-100 text-sm">For Survey: {survey.name}</p>
        </div>
        
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          <div>
            <label htmlFor="context-type" className="block text-sm font-medium text-gray-700">
              Select Context Type
            </label>
            <p className="mt-1 text-sm text-gray-500">
                Choose the context that best fits this survey to help the AI generate more relevant insights.
            </p>
            <select
              id="context-type"
              name="context-type"
              className="mt-2 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
              value={contextType}
              onChange={(e) => setContextType(e.target.value)}
            >
              {CONTEXT_TYPES.map((type) => (
                <option key={type.value} value={type.value}>
                  {type.label}
                </option>
              ))}
            </select>
          </div>

          <div className="flex justify-end space-x-3">
             <Link
                href="/admin"
                className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Cancel
              </Link>
            <button
              type="submit"
              disabled={loading}
              className={`inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              {loading ? 'Starting...' : 'Generate'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
