import { error } from "@sveltejs/kit";
import type { ComputerLoginRow, ComputerRow } from "$lib/computers";
import { requireComputersAccess } from "$lib/computersLoad";
import { pb } from "$lib/pocketbase";
import type { PageLoad } from "./$types";

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
  // Every sign-in of the computer: they are deleted after 30 days.
  const logins = (await pb.send("/api/computers/logins", {
    method: "GET",
    query: { computer: params.id },
  })) as ComputerLoginRow[];
  return { computer, logins };
};
