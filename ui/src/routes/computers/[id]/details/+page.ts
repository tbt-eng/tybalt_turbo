import { error } from "@sveltejs/kit";
import type { ComputerLoginRow, ComputerPage, ComputerRow } from "$lib/computers";
import { requireComputersAccess } from "$lib/computersLoad";
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
  // The server names the matched user, which a holder of the computers claim
  // can't read through the API.
  const logins = (await pb.send("/api/computers/logins", {
    method: "GET",
    query: { computer: params.id, perPage: LOGINS_SHOWN },
  })) as ComputerPage<ComputerLoginRow>;
  return { computer, logins: logins.items };
};
