'use client';

import { useEffect, useState, use } from 'react';
import { Insight } from '@/types';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function InsightDetails({ params }: { params: Promise<{ id: string }> }) {
  const router = useRouter();
  const { id } = use(params);
  const [insight, setInsight] = useState<Insight | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`/next-api/admin/insights/${id}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.data) setInsight(data.data);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="text-xl font-bold tracking-tight uppercase">Loading_Data...</div>;
  if (!insight) return <div>Insight not found</div>;

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6 flex justify-between items-center">
        <div>
          <h3 className="text-lg leading-6 font-medium text-gray-900">Insight Report</h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">
            Status: {insight.status} | Context: {insight.context_type}
          </p>
        </div>
        <div className="flex space-x-4">
            <button 
                type="button"
                onClick={() => router.back()}
                className="inline-flex items-center px-4 py-2 border-2 border-black text-sm font-bold uppercase text-black bg-white hover:bg-gray-100 transition-colors"
            >
                Back
            </button>
        </div>
      </div>
      <div className="border-t border-gray-200 px-4 py-5 sm:p-0">
        <dl className="sm:divide-y sm:divide-gray-200">
          <div className="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Overall Summary</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2 whitespace-pre-wrap">
              {insight.analysis || 'No summary generated yet.'}
            </dd>
          </div>
          
          <div className="py-4 sm:py-5 sm:px-6">
            <h4 className="text-md font-bold text-gray-900 mb-4">Question Analysis</h4>
            <div className="space-y-6">
                {insight.batches?.map((batch, index) => (
                    <div key={index} className="bg-gray-50 p-4 rounded-lg">
                        <p className="font-semibold text-gray-700 mb-2">Q: {batch.question.text}</p>
                        <p className="text-sm text-gray-600 mb-2 italic">Batch {batch.batch_number}</p>
                         <p className="text-gray-800 whitespace-pre-wrap">{batch.summary}</p>
                    </div>
                ))}
            </div>
          </div>
        </dl>
      </div>
    </div>
  );
}
