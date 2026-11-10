# Absorb system

The parts make up the UI portion of the absorb system.

1. A full-view component (AbsorbEditor.svelte) that shows the target record (for
example a clients record) at the top using a snippet described in the containing
+page.svelte file and a heading just below that says "records to absorb". Under
this heading is a list of zero or more records to be absorbed. Initially there
are zero items in the list. There is a search box (using the DSAutoComplete.svelte
component when the page passes a search index, otherwise the DSSelector.svelte
component) to immediately beneath the heading to find records from the same
collection as the target record (for example a clients record) and pressing
enter will add the selected record to the list of records to be absorbed. There
is a button at the bottom to absorb and another to cancel. Clicking the absorb
button absorbs the records and opens the Absorb Actions page (/absorb/actions),
where the pending absorb action is committed or undone. A contact merge opened
from a client's page stays on the page and shows the pending action there. While
a collection has a pending absorb action, the component shows that action in
place of the editor.

2. A +page.svelte file for each collection that supports absorption (clients, vendors and client contacts), contained within the /absorb directory under the collection's id directory (i.e. \[cid\]/absorb/+page.svelte under the clients route). Only the client contacts route also has a +page.ts file. The AbsorbEditor.svelte component is used in the +page.svelte file. Arguments are passed to the AbsorbEditor.svelte component, namely the collectionName, the targetRecordId (which is the id of the absorbing record pulled from the URL), the availableRecords and an optional autoCompleteIndex. A snippet is also provided to the AbsorbEditor.svelte component to render the record to be absorbed. The snippet is defined in the +page.svelte file on the route.

3. A button to absorb the records which appears on the regular list view for
clients and vendors, and beside each contact on the client's page. Options for
the icon are:

    - fluent-emoji-high-contrast:sponge
    - mdi:merge
    - material-symbols:merge

The buttons use the mdi:merge icon.
