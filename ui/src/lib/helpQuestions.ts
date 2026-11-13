// The questions that the home page shows as examples of what to ask Help. Each
// one is for every person, so none is about a task that needs a claim.
export const helpQuestions = [
  "How do I submit a time sheet?",
  "How do I enter my time?",
  "How do I change a time sheet that I submitted?",
  "How do I see my time off?",
  "How do I choose my manager?",
  "How do I create a purchase order?",
  "When do I need a purchase order?",
  "What are the types of purchase order?",
  "How do I enter an expense?",
  "How do I change an expense?",
  "How do I start a new project?",
  "What happens after I submit an expense?",
];

export function randomHelpQuestion(): string {
  return helpQuestions[Math.floor(Math.random() * helpQuestions.length)];
}
