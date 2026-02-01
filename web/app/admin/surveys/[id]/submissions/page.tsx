'use client';

import { useEffect, useState, use } from 'react';
import { Survey, Submission } from '@/types';
import Link from 'next/link';

export default function SurveySubmissions({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [submissions, setSubmissions] = useState<Submission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
        try {
            const [surveyRes, submissionsRes] = await Promise.all([
                fetch(`/next-api/admin/surveys/${id}`),
                fetch(`/next-api/admin/submissions?surveyId=${id}`)
            ]);
            
            if (surveyRes.ok) {
                const data = await surveyRes.json();
                setSurvey(data.data);
            }
            if (submissionsRes.ok) {
                const data = await submissionsRes.json();
                setSubmissions(data.data || []);
            }
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    }
    fetchData();
  }, [id]);

  if (loading) return <div className="text-xl font-bold tracking-tight uppercase">Loading_Data...</div>;
  if (!survey) return <div>Survey not found</div>;

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6 flex justify-between items-center">
        <div>
          <h3 className="text-lg leading-6 font-medium text-gray-900">Submissions for {survey.name}</h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">
            Total Submissions: {submissions.length}
          </p>
        </div>
        <Link href="/admin" className="text-indigo-600 hover:text-indigo-900">
            Back to Dashboard
        </Link>
      </div>
      <div className="flex flex-col">
          <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
              <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Submission Date
                      </th>
                      {survey.questions.map(q => (
                          <th key={q.id} scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              {q.text}
                          </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {submissions.map((submission) => (
                      <tr key={submission.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                           {new Date(submission.created_at).toLocaleString()}
                        </td>
                        {survey.questions.map(q => {
                            const response = submission.responses.find(r => r.question_id === q.id);
                            return (
                                <td key={q.id} className="px-6 py-4 whitespace-normal text-sm text-gray-900">
                                    {response?.answer || '-'}
                                </td>
                            );
                        })}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
      </div>
    </div>
  );
}
