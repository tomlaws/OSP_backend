'use client';

import { useEffect, useState, Suspense } from 'react';
import Link from 'next/link';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';
import { Survey, Insight } from '@/types';

export default function AdminDashboard() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <AdminDashboardContent />
    </Suspense>
  );
}

function AdminDashboardContent() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  
  const page = Number(searchParams.get('page')) || 1;
  const [surveys, setSurveys] = useState<Survey[]>([]);
  const [insights, setInsights] = useState<Insight[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const pageSize = 5;

  const updatePage = (newPage: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('page', newPage.toString());
    router.push(`${pathname}?${params.toString()}`);
  };

  useEffect(() => {
    fetchData();
  }, [page]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const offset = (page - 1) * pageSize;
      const [surveysRes, insightsRes] = await Promise.all([
        fetch(`/next-api/admin/surveys?offset=${offset}&limit=${pageSize}`),
        fetch('/next-api/admin/insights'),
      ]);

      if (surveysRes.ok) {
        const data = await surveysRes.json();
        setSurveys(data.data || []);
        setTotal(data.total || 0);
      }
      if (insightsRes.ok) {
        const data = await insightsRes.json();
        setInsights(data.data || []);
      }
    } catch (error) {
      console.error('Failed to fetch data', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteSurvey = async (id: string) => {
    if (!confirm('Are you sure?')) return;
    try {
      await fetch(`/next-api/admin/surveys/${id}`, { method: 'DELETE' });
      fetchData();
    } catch (e) {
      console.error(e);
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <div className="md:flex md:items-center md:justify-between mb-6">
        <div className="flex-1 min-w-0">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Surveys
          </h2>
        </div>
        <div className="mt-4 flex md:mt-0 md:ml-4">
          <Link
            href="/admin/create"
            className="ml-3 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Create New Survey
          </Link>
        </div>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {surveys.map((survey) => {
            const insight = insights.find(i => i.survey_id === survey.id);
            const shareUrl = `${window.location.origin}/survey/${survey.token}`;

            return (
              <li key={survey.id}>
                <div className="px-4 py-4 sm:px-6">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-medium text-indigo-600 truncate">
                      {survey.name}
                    </p>
                    <div className="ml-2 flex-shrink-0 flex">
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        Active
                      </span>
                    </div>
                  </div>
                  <div className="mt-2 sm:flex sm:justify-between">
                    <div className="sm:flex">
                      <p className="flex items-center text-sm text-gray-500">
                        Token: {survey.token}
                      </p>
                    </div>
                  </div>
                   <div className="mt-4 flex space-x-4 text-sm">
                      <a href={shareUrl} target="_blank" className="text-indigo-600 hover:text-indigo-900">
                        Share Link
                      </a>
                      <button 
                        onClick={() => handleDeleteSurvey(survey.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Delete
                      </button>

                      <Link href={`/admin/surveys/${survey.id}/submissions`} className="text-gray-600 hover:text-gray-900 font-medium">
                        View Submissions
                      </Link>
                      
                      {insight ? (
                        <Link href={`/admin/insights/${insight.id}`} className="text-purple-600 hover:text-purple-900 font-medium">
                          View Insights ({insight.status})
                        </Link>
                      ) : (
                        <Link 
                          href={`/admin/surveys/${survey.id}/generate-insight`}
                          className="text-blue-600 hover:text-blue-900"
                        >
                          Generate Insight
                        </Link>
                      )}
                      
                   </div>
                </div>
              </li>
            );
          })}
          {surveys.length === 0 && (
             <div className="px-4 py-4 sm:px-6 text-center text-gray-500">
                No surveys found. Create one to get started.
             </div>
          )}
        </ul>
      </div>

      {total > pageSize && (
        <div className="flex justify-between items-center mt-4 pb-10">
          <button
            onClick={() => updatePage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="px-4 py-2 border rounded-md bg-white disabled:opacity-50"
          >
            Previous
          </button>
          <span>
            Page {page} of {Math.ceil(total / pageSize)}
          </span>
          <button
            onClick={() => updatePage(page + 1)}
            disabled={page >= Math.ceil(total / pageSize)}
            className="px-4 py-2 border rounded-md bg-white disabled:opacity-50"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
