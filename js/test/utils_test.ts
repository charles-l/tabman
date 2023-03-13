import { assertEquals } from "https://deno.land/std@0.179.0/testing/asserts.ts";
import { humanTimeDiff } from "../../add-on/utils.js";

Deno.test("test humanTimeDiff", () => {
  assertEquals(humanTimeDiff(1), "1 second");
  assertEquals(humanTimeDiff(60), "1 minute");
  assertEquals(humanTimeDiff(60 * 60), "1 hour");
  assertEquals(humanTimeDiff(60 * 60 * 24), "1 day");
  assertEquals(humanTimeDiff(60 * 60 * 24 * 365), "365 days");
  assertEquals(humanTimeDiff(30), "30 seconds");
  assertEquals(humanTimeDiff(400), "6 minutes");
  assertEquals(humanTimeDiff(10000), "2 hours");
  assertEquals(humanTimeDiff(340600), "3 days");
  assertEquals(humanTimeDiff(345600), "4 days");
});
