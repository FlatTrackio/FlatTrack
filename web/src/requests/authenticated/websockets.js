import common from "@/common/common";

function GetWebSocket(group, id, ...searchParams) {
  let u = new URL(window.location.href);
  u.protocol = "wss:";
  if (window.location.protocol === "http:") {
    u.protocol = "ws:";
  }
  u.pathname = `/api/ws/${group}/${id}`;
  searchParams.forEach((p) => {
    u.searchParams.set(p.name, p.value);
  });
  return new WebSocket(u.toString());
}

export default {
  GetWebSocket,
};
