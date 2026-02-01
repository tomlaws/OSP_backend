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
      <div className="py-20 border-2 border-black p-8 bg-gray-50">
        <h2 className="text-4xl font-black text-black mb-4 uppercase tracking-tighter">Transmission Received</h2>
        <div className="w-16 h-2 bg-[#D80000] mb-6"></div>
        <p className="text-xl font-bold text-gray-800">Your response has been recorded in the database.</p>
        <p className="text-sm font-mono text-gray-500 mt-4 uppercase">Session Terminated.</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-16">
      {survey.questions.map((q, idx) => (
        <div key={q.id} className="space-y-4">
          <div className="flex items-baseline gap-4">
            <span className="text-xs font-mono font-bold text-gray-400">0{idx + 1}.</span>
            <label className="block text-2xl font-bold text-black uppercase leading-tight tracking-tight">
                {q.text}
            </label>
          </div>
          
          <div className="ml-0 md:ml-10">
          {q.type === 'TEXTBOX' ? (
            <textarea
              required
              maxLength={q.specification?.max_length || 100}
              className="mt-1 block w-full bg-transparent border-b-2 border-gray-300 focus:border-black focus:ring-0 p-2 text-xl font-medium outline-none transition-colors rounded-none placeholder-gray-300"
              rows={2}
              placeholder="Type your answer here..."
              onChange={(e) => setAnswers({ ...answers, [q.id!]: e.target.value })}
            />
          ) : q.type === 'LIKERT' ? (
            <div className="pt-2">
              <div className="flex flex-col md:flex-row items-center justify-between gap-4 border-2 border-gray-100 p-6 bg-gray-50">
                {q.specification?.min_label && (
                    <span className="text-xs font-bold uppercase tracking-widest text-gray-500 flex-1 text-right md:text-left">{q.specification.min_label}</span>
                )}
                
                <div className="flex space-x-4 md:space-x-8">
                    {Array.from(
                    { length: (q.specification?.max || 5) - (q.specification?.min || 1) + 1 }, 
                    (_, i) => (q.specification?.min || 1) + i
                    ).map((val) => (
                    <label key={val} className="cursor-pointer group flex flex-col items-center gap-2">
                        <input
                            type="radio"
                            name={q.id}
                            required
                            value={val}
                            onChange={(e) => setAnswers({ ...answers, [q.id!]: e.target.value })}
                            className="peer sr-only"
                        />
                        <div className="w-10 h-10 md:w-12 md:h-12 border-2 border-gray-300 flex items-center justify-center text-lg font-bold text-gray-400 peer-checked:border-black peer-checked:bg-black peer-checked:text-white transition-all group-hover:border-gray-500">
                            {val}
                        </div>
                    </label>
                    ))}
                </div>

                {q.specification?.max_label && (
                    <span className="text-xs font-bold uppercase tracking-widest text-gray-500 flex-1 text-left md:text-right">{q.specification.max_label}</span>
                )}
              </div>
            </div>
          ) : q.type === 'MULTIPLE_CHOICE' ? (
            <div className="space-y-3 pt-2">
              {q.specification?.options?.map((opt) => (
                <label key={opt} className="flex items-center group cursor-pointer p-3 border border-transparent hover:border-gray-200 transition-colors">
                  <input
                    type="radio"
                    name={q.id}
                    required
                    value={opt}
                    onChange={(e) => setAnswers({ ...answers, [q.id!]: e.target.value })}
                    className="peer sr-only"
                  />
                  <div className="w-5 h-5 border-2 border-gray-300 peer-checked:border-[#D80000] peer-checked:bg-[#D80000] mr-4 flex-shrink-0 transition-colors"></div>
                  <span className="text-lg font-medium text-gray-600 peer-checked:text-black peer-checked:font-bold">
                    {opt}
                  </span>
                </label>
              ))}
            </div>
          ) : null}
          </div>
        </div>
      ))}

      <div className="pt-12 border-t-2 border-black flex justify-end">
        <button
          type="submit"
          disabled={submitting}
          className="inline-flex justify-center py-4 px-12 border-2 border-black text-lg font-bold uppercase tracking-widest text-white bg-black hover:bg-[#D80000] hover:border-[#D80000] hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,0.2)] transition-all disabled:opacity-50"
        >
          {submitting ? 'Authenticating...' : 'Submit Data'}
        </button>
      </div>
    </form>
  );
}
