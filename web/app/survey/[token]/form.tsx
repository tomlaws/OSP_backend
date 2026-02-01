'use client';

import { useState } from 'react';
import { Survey, CreateSubmissionRequest, SubmissionResponse } from '@/types';

export default function SurveyForm({ survey }: { survey: Survey }) {
  const [answers, setAnswers] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    const responses: SubmissionResponse[] = Object.entries(answers).map(
      ([qID, ans]) => ({
        question_id: qID,
        answer: ans,
      })
    );

    const payload: CreateSubmissionRequest = {
      survey_token: survey.token,
      responses,
    };

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/submissions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        setSubmitted(true);
      } else {
        alert('Submission failed. Please try again.');
      }
    } catch (error) {
      console.error(error);
      alert('Network error');
    } finally {
      setSubmitting(false);
    }
  };

  if (submitted) {
    return (
      <div className="text-center py-10">
        <h2 className="text-2xl font-bold text-green-600 mb-4">Thank You!</h2>
        <p className="text-gray-600">Your response has been recorded.</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      {survey.questions.map((q) => (
        <div key={q.id} className="space-y-2">
          <label className="block text-lg font-medium text-gray-900">
            {q.text}
          </label>
          
          {q.type === 'TEXTBOX' ? (
            <textarea
              required
              maxLength={q.specification?.max_length || 100}
              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-3 focus:ring-indigo-500 focus:border-indigo-500"
              rows={3}
              onChange={(e) => setAnswers({ ...answers, [q.id!]: e.target.value })}
            />
          ) : q.type === 'LIKERT' ? (
            <div className="flex items-center space-x-4">
              {q.specification?.min_label && (
                <span className="text-sm font-medium text-gray-500">{q.specification.min_label}</span>
              )}
              <div className="flex space-x-6">
                {Array.from(
                  { length: (q.specification?.max || 5) - (q.specification?.min || 1) + 1 }, 
                  (_, i) => (q.specification?.min || 1) + i
                ).map((val) => (
                  <div key={val} className="flex flex-col items-center">
                      <input
                        type="radio"
                        name={q.id}
                        required
                        value={val}
                        onChange={(e) => setAnswers({ ...answers, [q.id!]: e.target.value })}
                        className="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300"
                      />
                      <label className="mt-1 text-sm text-gray-700">{val}</label>
                  </div>
                ))}
              </div>
              {q.specification?.max_label && (
                <span className="text-sm font-medium text-gray-500">{q.specification.max_label}</span>
              )}
            </div>
          ) : q.type === 'MULTIPLE_CHOICE' ? (
            <div className="space-y-2">
              {q.specification?.options?.map((opt) => (
                <div key={opt} className="flex items-center">
                  <input
                    type="radio"
                    name={q.id}
                    required
                    value={opt}
                    onChange={(e) => setAnswers({ ...answers, [q.id!]: e.target.value })}
                    className="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300"
                  />
                  <label className="ml-3 block text-sm font-medium text-gray-700">
                    {opt}
                  </label>
                </div>
              ))}
            </div>
          ) : null}
        </div>
      ))}

      <div className="pt-5">
        <button
          type="submit"
          disabled={submitting}
          className="w-full inline-flex justify-center py-3 px-6 border border-transparent shadow-sm text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {submitting ? 'Submitting...' : 'Submit Responses'}
        </button>
      </div>
    </form>
  );
}
