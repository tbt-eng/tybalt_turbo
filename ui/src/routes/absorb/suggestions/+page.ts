import type { PageLoad } from "./$types";
import {
  absorbSuggestionLevels,
  isAbsorbSuggestionCollection,
  type AbsorbSuggestionCollection,
  type AbsorbSuggestionLevel,
} from "$lib/absorbSuggestions";

// The tab and the level are in the URL, so a reload keeps them:
// /absorb/suggestions?tab=vendors&level=loose. The page loads the buckets
// itself, so changing the tab keeps the work done in the other one.
export const load: PageLoad = ({ url }) => {
  const tabParam = url.searchParams.get("tab");
  const levelParam = url.searchParams.get("level");
  const tab: AbsorbSuggestionCollection = isAbsorbSuggestionCollection(tabParam)
    ? tabParam
    : "clients";
  const level: AbsorbSuggestionLevel =
    absorbSuggestionLevels.find((option) => option.id === levelParam)?.id ?? "normal";
  return { tab, level };
};
