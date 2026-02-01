'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { CreateSurveyRequest, Question } from '@/types';

export default function CreateSurvey() {
  const router = useRouter();
  const [name, setName] = useState('');
  const [questions, setQuestions] = useState<Question[]>([]);

  const addQuestion = () => {
    setQuestions([
      ...questions,
      {
        type: 'TEXTBOX',
        text: '',
        specification: { max_length: 100 },
      },
    ]);
  };

  const updateQuestion = (index: number, field: keyof Question, value: any) => {
    const newQuestions = [...questions];
    newQuestions[index] = { ...newQuestions[index], [field]: value };
    // Set defaults when type changes
    if (field === 'type') {
       if (value === 'TEXTBOX') {
           newQuestions[index].specification = { max_length: 100 };
       } else if (value === 'LIKERT') {
           newQuestions[index].specification = { min: 1, max: 5 };
       } else {
           newQuestions[index].specification = {};
       }
    }
    setQuestions(newQuestions);
  };

  const updateSpec = (index: number, field: string, value: any) => {
    const newQuestions = [...questions];
    newQuestions[index].specification = {
        ...newQuestions[index].specification,
        [field]: value
    };
    // Ensure min is set if not present when setting max
    if (field === 'max' && !newQuestions[index].specification.min) {
        newQuestions[index].specification.min = 1;
    }
    setQuestions(newQuestions);
  };
  
  const removeQuestion = (index: number) => {
      setQuestions(questions.filter((_, i) => i !== index));
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const payload: CreateSurveyRequest = {
      name,
      questions,
    };

    try {
      const res = await fetch('/next-api/admin/surveys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        router.push('/admin');
      } else {
        alert('Failed to create survey');
      }
    } catch (error) {
      console.error(error);
      alert('Error creating survey');
    }
  };

  return (
    <div className="max-w-4xl mx-auto pb-24">
      <div className="mb-12 border-b-4 border-black pb-4">
        <h1 className="text-4xl font-black uppercase tracking-tighter">New Survey Construction</h1>
      </div>
      
      <form onSubmit={handleSubmit} className="space-y-12">
        <div>
          <label className="block text-xs font-bold uppercase tracking-widest text-black mb-2">Survey Designation</label>
          <input
            type="text"
            required
            className="block w-full border-2 border-black bg-white p-4 text-xl font-bold focus:ring-0 focus:border-[#D80000] outline-none transition-colors"
            placeholder="ENTER SURVEY NAME"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="space-y-8">
          <div className="flex items-center justify-between border-b-2 border-black pb-2">
            <h3 className="text-2xl font-bold uppercase tracking-tight">Questions Sequence</h3>
            <span className="font-mono text-sm bg-black text-white px-2 py-1">{questions.length} ITEMS</span>
          </div>

          {questions.map((q, idx) => (
            <div key={idx} className="border-2 border-black p-6 relative bg-white group hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] transition-shadow duration-200">
                <div className="absolute -top-4 -left-4 bg-black text-white w-8 h-8 flex items-center justify-center font-bold font-mono border-2 border-white">
                    {idx + 1}
                </div>
                
                <button 
                  type="button" 
                  onClick={() => removeQuestion(idx)} 
                  className="absolute top-4 right-4 text-xs font-bold uppercase hover:text-[#D80000] hover:underline"
                >
                    [Remove]
                </button>
                
              <div className="grid grid-cols-1 md:grid-cols-12 gap-6 mt-2">
                <div className="col-span-12 md:col-span-4">
                  <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Type</label>
                  <select
                    className="block w-full border-2 border-black bg-white p-3 font-bold focus:ring-0 focus:border-[#D80000] outline-none appearance-none rounded-none"
                    value={q.type}
                    onChange={(e) => updateQuestion(idx, 'type', e.target.value)}
                  >
                    <option value="TEXTBOX">TEXTBOX</option>
                    <option value="MULTIPLE_CHOICE">MULTIPLE CHOICE</option>
                    <option value="LIKERT">LIKERT SCALE</option>
                  </select>
                </div>
                
                <div className="col-span-12 md:col-span-8">
                  <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">Prompt</label>
                  <input
                    type="text"
                    required
                    className="block w-full border-b-2 border-black bg-transparent p-3 font-medium focus:ring-0 focus:border-[#D80000] outline-none rounded-none"
                    placeholder="Enter question text here..."
                    value={q.text}
                    onChange={(e) => updateQuestion(idx, 'text', e.target.value)}
                  />
                </div>

                {q.type === 'TEXTBOX' && (
                  <div className="col-span-12 md:col-span-6 bg-gray-50 p-4 border border-black">
                    <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                      Max Character Length
                    </label>
                    <input
                      type="number"
                      className="block w-full bg-white border border-gray-400 p-2 font-mono text-sm"
                      value={q.specification?.max_length || 100}
                      onChange={(e) =>
                        updateSpec(idx, 'max_length', parseInt(e.target.value))
                      }
                    />
                  </div>
                )}

                {q.type === 'LIKERT' && (
                  <div className="col-span-12 bg-gray-50 p-4 border border-black">
                     <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div>
                        <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                            Min Value
                        </label>
                        <input
                            type="number"
                            className="block w-full bg-white border border-gray-400 p-2 font-mono text-sm"
                            value={q.specification?.min || 1}
                            onChange={(e) =>
                            updateSpec(idx, 'min', parseInt(e.target.value))
                            }
                        />
                        </div>
                        <div>
                        <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                            Max Value
                        </label>
                        <input
                            type="number"
                            className="block w-full bg-white border border-gray-400 p-2 font-mono text-sm"
                            value={q.specification?.max || 5}
                            min={(q.specification?.min || 1) + 1}
                            max={10}
                            onChange={(e) =>
                            updateSpec(idx, 'max', parseInt(e.target.value))
                            }
                        />
                        </div>
                        <div>
                        <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                            Min Label
                        </label>
                        <input
                            type="text"
                            className="block w-full bg-white border border-gray-400 p-2 font-mono text-sm"
                            placeholder="e.g. Strongly Disagree"
                            value={q.specification?.min_label || ''}
                            onChange={(e) =>
                            updateSpec(idx, 'min_label', e.target.value)
                            }
                        />
                        </div>
                        <div>
                        <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                            Max Label
                        </label>
                        <input
                            type="text"
                            className="block w-full bg-white border border-gray-400 p-2 font-mono text-sm"
                            placeholder="e.g. Strongly Agree"
                            value={q.specification?.max_label || ''}
                            onChange={(e) =>
                            updateSpec(idx, 'max_label', e.target.value)
                            }
                        />
                        </div>
                    </div>
                  </div>
                )}

                {q.type === 'MULTIPLE_CHOICE' && (
                  <div className="col-span-12 bg-gray-50 p-4 border border-black">
                    <label className="block text-xs font-bold uppercase tracking-widest text-gray-500 mb-2">
                      Options (comma separated)
                    </label>
                    <input
                      type="text"
                      className="block w-full bg-white border border-gray-400 p-2 font-mono text-sm"
                      placeholder="Option 1, Option 2, Option 3"
                      onChange={(e) =>
                        updateSpec(idx, 'options', e.target.value.split(',').map((s) => s.trim()))
                      }
                    />
                  </div>
                )}
              </div>
            </div>
          ))}

          <button
            type="button"
            onClick={addQuestion}
            className="w-full py-4 border-2 border-dashed border-gray-400 text-sm font-bold uppercase text-gray-500 hover:border-black hover:text-black hover:bg-gray-50 transition-all"
          >
            + Add Another Question
          </button>
        </div>

        <div className="flex justify-end pt-8 border-t-2 border-black">
          <button
            type="submit"
            className="inline-flex justify-center py-4 px-8 border-2 border-black text-sm font-bold uppercase text-white bg-black hover:bg-[#D80000] hover:border-[#D80000] shadow-[4px_4px_0px_0px_rgba(0,0,0,0.5)] hover:shadow-none hover:translate-x-[2px] hover:translate-y-[2px] transition-all"
          >
            Create Survey
          </button>
        </div>
      </form>
    </div>
  );
}
