import { error } from "@sveltejs/kit";
import type { ComputerRow } from "$lib/computers";
import { requireComputersAccess } from "$lib/computersLoad";
import type { ComputerLoginsResponse } from "$lib/pocketbase-types";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";

// The newest sign-ins that the page shows.
const LOGINS_SHOWN = 20;

export const load: PageLoad = async ({ params }) => {
  await requireComputersAccess();
  let computer: ComputerRow;
  try {
    computer = (await pb.send(`/api/computers/${params.id}`, { method: "GET" })) as ComputerRow;
  } catch (err) {
    console.error(`loading computer details: ${err}`);
    error(404, "Computer not found.");
  }
  const logins = await pb
    .collection("computer_logins")
    .getList<ComputerLoginsResponse>(1, LOGINS_SHOWN, {
      filter: pb.filter("computer = {:computer}", { computer: params.id }),
      sort: "-created",
    });
  return { computer, logins: logins.items };
};
