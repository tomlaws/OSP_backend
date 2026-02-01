'use client';

import { useEffect, useState, Suspense } from 'react';
import Link from 'next/link';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';
import { Survey, Insight } from '@/types';

export default function AdminDashboard() {
  return (
    <Suspense fallback={<div className="text-xl font-bold tracking-tight uppercase">Loading_Data...</div>}>
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

  if (loading) return <div className="text-xl font-bold tracking-tight uppercase">Loading_Data...</div>;

  return (
    <div>
      <div className="md:flex md:items-center md:justify-between mb-12">
        <div className="flex-1 min-w-0 border-b-4 border-black pb-2">
          <h2 className="text-4xl font-black uppercase tracking-tighter text-black">
            Surveys
          </h2>
        </div>
        <div className="mt-4 flex md:mt-0 md:ml-12 self-end">
          <Link
            href="/admin/create"
            className="inline-flex items-center px-6 py-3 border-2 border-black text-sm font-bold uppercase text-white bg-black hover:bg-[#D80000] hover:border-[#D80000] transition-colors duration-200"
          >
            + Create New
          </Link>
        </div>
      </div>

      <div className="border-t-2 border-black">
        <ul className="divide-y-2 divide-black">
          {surveys.map((survey) => {
            const insight = insights.find(i => i.survey_id === survey.id);
            const shareUrl = `${window.location.origin}/survey/${survey.token}`;

            return (
              <li key={survey.id} className="group hover:bg-gray-50 transition-colors duration-200">
                <div className="py-6">
                  <div className="flex items-start justify-between">
                    <div>
                        <p className="text-2xl font-bold uppercase tracking-tight text-black mb-1">
                        {survey.name}
                        </p>
                        <p className="text-sm font-mono text-gray-500 uppercase tracking-widest">
                            Token: {survey.token}
                        </p>
                    </div>
                    
                    <div className="ml-2 flex-shrink-0 flex">
                      <span className="px-3 py-1 text-xs font-bold uppercase border border-black bg-white text-black">
                        Active
                      </span>
                    </div>
                  </div>
                  
                   <div className="mt-6 flex flex-wrap gap-x-8 gap-y-4 text-sm font-bold uppercase tracking-wide">
                      <a href={shareUrl} target="_blank" className="text-black hover:text-[#D80000] hover:underline underline-offset-4">
                        Share Link
                      </a>
                      
                      <Link href={`/admin/surveys/${survey.id}/submissions`} className="text-black hover:text-[#D80000] hover:underline underline-offset-4">
                        View Submissions
                      </Link>
                      
                      {insight ? (
                        <Link href={`/admin/insights/${insight.id}`} className="text-black hover:text-[#D80000] hover:underline underline-offset-4">
                          Insights <span className="text-[#D80000]">({insight.status})</span>
                        </Link>
                      ) : (
                        <Link 
                          href={`/admin/surveys/${survey.id}/generate-insight`}
                          className="text-black hover:text-[#D80000] hover:underline underline-offset-4"
                        >
                          Generate Insight
                        </Link>
                      )}

                      <button 
                        onClick={() => handleDeleteSurvey(survey.id)}
                        className="text-gray-400 hover:text-[#D80000] ml-auto"
                      >
                        Delete
                      </button>
                      
                   </div>
                </div>
              </li>
            );
          })}
          {surveys.length === 0 && (
             <div className="py-12 text-center text-gray-500 font-bold uppercase">
                No surveys found. Create one to get started.
             </div>
          )}
        </ul>
      </div>

      {total > pageSize && (
        <div className="flex justify-between items-center mt-12 py-6 border-t-2 border-black">
          <button
            onClick={() => updatePage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="px-6 py-3 border-2 border-black text-sm font-bold uppercase disabled:opacity-30 hover:bg-black hover:text-white transition-colors"
          >
            Previous
          </button>
          <span className="font-mono text-sm">
            PAGE {page} / {Math.ceil(total / pageSize)}
          </span>
          <button
            onClick={() => updatePage(page + 1)}
            disabled={page >= Math.ceil(total / pageSize)}
            className="px-6 py-3 border-2 border-black text-sm font-bold uppercase disabled:opacity-30 hover:bg-black hover:text-white transition-colors"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
