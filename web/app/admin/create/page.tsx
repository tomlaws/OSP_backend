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
    <div className="max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Create New Survey</h1>
      <form onSubmit={handleSubmit} className="space-y-6 bg-white p-6 shadow rounded-lg">
        <div>
          <label className="block text-sm font-medium text-gray-700">Survey Name</label>
          <input
            type="text"
            required
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="space-y-4">
          <h3 className="text-lg font-medium">Questions</h3>
          {questions.map((q, idx) => (
            <div key={idx} className="border p-4 rounded-md relative bg-gray-50">
                <button 
                  type="button" 
                  onClick={() => removeQuestion(idx)} 
                  className="absolute top-2 right-2 text-red-500 text-sm"
                >
                    Remove
                </button>
              <div className="grid grid-cols-1 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Question Type</label>
                  <select
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                    value={q.type}
                    onChange={(e) => updateQuestion(idx, 'type', e.target.value)}
                  >
                    <option value="TEXTBOX">Textbox</option>
                    <option value="MULTIPLE_CHOICE">Multiple Choice</option>
                    <option value="LIKERT">Likert Scale</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Question Text</label>
                  <input
                    type="text"
                    required
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                    value={q.text}
                    onChange={(e) => updateQuestion(idx, 'text', e.target.value)}
                  />
                </div>

                {q.type === 'TEXTBOX' && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Max Length
                    </label>
                    <input
                      type="number"
                      className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                      value={q.specification?.max_length || 100}
                      onChange={(e) =>
                        updateSpec(idx, 'max_length', parseInt(e.target.value))
                      }
                    />
                  </div>
                )}

                {q.type === 'LIKERT' && (
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Min Value
                      </label>
                      <input
                        type="number"
                        className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                        value={q.specification?.min || 1}
                        onChange={(e) =>
                          updateSpec(idx, 'min', parseInt(e.target.value))
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Max Value
                      </label>
                      <input
                        type="number"
                        className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                        value={q.specification?.max || 5}
                        min={(q.specification?.min || 1) + 1}
                        max={10}
                        onChange={(e) =>
                          updateSpec(idx, 'max', parseInt(e.target.value))
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Min Label (Optional)
                      </label>
                      <input
                        type="text"
                        className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                        placeholder="e.g. Strongly Disagree"
                        value={q.specification?.min_label || ''}
                        onChange={(e) =>
                          updateSpec(idx, 'min_label', e.target.value)
                        }
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Max Label (Optional)
                      </label>
                      <input
                        type="text"
                        className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
                        placeholder="e.g. Strongly Agree"
                        value={q.specification?.max_label || ''}
                        onChange={(e) =>
                          updateSpec(idx, 'max_label', e.target.value)
                        }
                      />
                    </div>
                  </div>
                )}

                {q.type === 'MULTIPLE_CHOICE' && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Options (comma separated)
                    </label>
                    <input
                      type="text"
                      className="mt-1 block w-full border-gray-300 rounded-md shadow-sm p-2"
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
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-indigo-700 bg-indigo-100 hover:bg-indigo-200"
          >
            Add Question
          </button>
        </div>

        <div className="flex justify-end pt-4">
          <button
            type="submit"
            className="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
          >
            Create Survey
          </button>
        </div>
      </form>
    </div>
  );
}
