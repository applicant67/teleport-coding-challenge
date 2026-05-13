# Historical submissions — master index

This page lists every coding-challenge submission we have captured, across all tracks. It is the single entry point into the corpus.

## How to read this table

- **AI score**: average of two independent grader passes against the rubric in [`SCORING_RUBRIC.md`](SCORING_RUBRIC.md). Range −22 to +11. The two raw scores are shown small after the average. ⚠️ flags rows where the graders disagreed by more than 2 points — treat the average with extra skepticism on those rows. ⚙ flags recovered (gh-archive-only) repos where the grader saw less than full source. ◈ flags repos whose source was recovered from the Software Heritage Archive after deletion from GitHub — graded on full source, just sourced via SWH rather than the original GitHub.
- **Hired**: ✓ means the submitter is a *current Teleport reviewer*, i.e. they later joined the company and now review challenge submissions themselves. Absence does not mean rejection — most public submissions are from people whose hire outcome we cannot verify.
- **Recovered**: rows under the recovered heading per track are submissions whose repos are no longer public on GitHub; we reconstructed activity from `data.gharchive.org`. See each repo's `RECONSTRUCTED.md`.
- **Level**: where the submission was scoped explicitly (L1–L5 in `tracks/systems/`), it is shown. Blank means the level was not stated in the captured artifact.
- **Date range**: first and last captured event timestamps. For recovered repos this is the gh-archive window only.
- **PR comments**: count of captured PR review comments — a rough proxy for how much reviewer signal the row carries.

See [`SCORING_RUBRIC.md`](SCORING_RUBRIC.md) for what the score means, the 11 published criteria, the 9 reviewer-feedback themes, and the methodology caveats (read these before drawing conclusions from a single score).

## Systems (job-worker / backend) (87 submissions)

| Repo | Date range | Level | Hired | PR cmts | AI score |
|---|---|---|---|---|---|
| [`ilyazz/jobs`](tracks/systems/corpus/raw/ilyazz__jobs/) | 2022-10-17 → 2022-10-28 | L5 |  | 89 | **10** <sub>A:10 B:10</sub> |
| [`kkloberdanz/teleport-challenge`](tracks/systems/corpus/raw/kkloberdanz__teleport-challenge/) | 2026-02-16 | L5 |  | 114 | **9.5** <sub>A:9 B:10</sub> |
| [`adalton/teleport-exercise`](tracks/systems/corpus/raw/adalton__teleport-exercise/) | 2021-12-21 | L5 |  | 158 | **8** ⚠️ <sub>A:10 B:6</sub> |
| [`GevorgGal/jobworker`](tracks/systems/corpus/raw/GevorgGal__jobworker/) | 2026-03-21 → 2026-04-02 | L4 |  | 146 | **8** <sub>A:8 B:8</sub> |
| [`mikewurtz/taskman`](tracks/systems/corpus/raw/mikewurtz__taskman/) | 2025-04-03 → 2025-04-15 | L5 |  | 154 | **8** <sub>A:8 B:8</sub> |
| [`xSudoNymx/job-manager`](tracks/systems/corpus/raw/xSudoNymx__job-manager/) | 2026-01-16 → 2026-01-29 | — |  | 90 | **8** <sub>A:8 B:8</sub> |
| [`rajansandeep/teleport-job-worker-svc`](tracks/systems/corpus/raw/rajansandeep__teleport-job-worker-svc/) | 2026-04-14 | L4 |  | 116 | **7.5** <sub>A:8 B:7</sub> |
| [`creack/telepilot`](tracks/systems/corpus/raw/creack__telepilot/) | 2024-10-07 → 2024-10-11 | — |  | 350 | **7** <sub>A:7 B:7</sub> |
| [`RichyHBM/teleport-challenge`](tracks/systems/corpus/raw/RichyHBM__teleport-challenge/) | 2025-04-17 → 2025-04-27 | L4 |  | 132 | **7** <sub>A:7 B:7</sub> |
| [`tjper/teleport`](tracks/systems/corpus/raw/tjper__teleport/) | 2022-03-16 | — |  | 143 | **7** ⚠️ <sub>A:9 B:5</sub> |
| [`GavinFrazar/teleport-challenge`](mirrors_extracted/systems/GavinFrazar__teleport-challenge/) | 2022-04-06 → 2022-04-16 | — | ✓ [†](#hired-note) | 0 ext | **6.5** ◈ <sub>A:7 B:6</sub> |
| [`MarkDHarris/teleport-systems-challenge`](tracks/systems/corpus/raw/MarkDHarris__teleport-systems-challenge/) | 2026-03-30 → 2026-04-03 | L4 |  | 108 | **6** ⚠️ <sub>A:8 B:4</sub> |
| [`MrChristianL/Teleport-Job-Worker-Service`](tracks/systems/corpus/raw/MrChristianL__Teleport-Job-Worker-Service/) | 2026-03-03 → 2026-03-10 | L4 |  | 199 | **6** <sub>A:6 B:6</sub> |
| [`Zephan92/teleport-handson`](tracks/systems/corpus/raw/Zephan92__teleport-handson/) | 2026-01-26 | — |  | 142 | **6** ⚠️ <sub>A:8 B:4</sub> |
| [`renatoaguimaraes/golang-job-scheduler`](tracks/systems/corpus/raw/renatoaguimaraes__golang-job-scheduler/) | 2021-04-25 → 2023-07-05 | L4 |  | 134 | **5.5** ⚠️ <sub>A:8 B:3</sub> |
| [`benmoss/job-worker-service`](tracks/systems/corpus/raw/benmoss__job-worker-service/) | 2026-01-09 | L4 |  | 71 | **5** <sub>A:4 B:6</sub> |
| [`nixpig/jobworker`](tracks/systems/corpus/raw/nixpig__jobworker/) | — | — |  | — | **5** ⚙ <sub>A:4 B:6</sub> |
| [`neildo/tjob`](tracks/systems/corpus/raw/neildo__tjob/) | 2024-09-11 → 2024-09-24 | — |  | 106 | **4.5** <sub>A:4 B:5</sub> |
| [`thompsy/go-linux-worker`](tracks/systems/corpus/raw/thompsy__go-linux-worker/) | 2021-04-23 → 2021-04-26 | — |  | — | **4.5** <sub>A:5 B:4</sub> |
| [`dustinspecker/linux-job-worker-service`](tracks/systems/corpus/raw/dustinspecker__linux-job-worker-service/) | — | — |  | — | **4** ⚙ <sub>A:4 B:4</sub> |
| [`zostay/hermione`](tracks/systems/corpus/raw/zostay__hermione/) | — | — |  | — | **4** ⚙ <sub>A:5 B:3</sub> |
| [`joshuarubin/teleport-job-worker`](tracks/systems/corpus/raw/joshuarubin__teleport-job-worker/) | 2024-09-10 → 2024-12-23 | L4 |  | 197 | **3.5** ⚠️ <sub>A:-1 B:8</sub> |
| [`supby/job-worker`](tracks/systems/corpus/raw/supby__job-worker/) | 2021-10-05 → 2021-10-06 | — |  | 51 | **3.5** <sub>A:4 B:3</sub> |
| [`trojal/mh-jobworker`](tracks/systems/corpus/raw/trojal__mh-jobworker/) | — | — |  | — | **3.5** ⚙ <sub>A:4 B:3</sub> |
| [`gstelang/job-worker-service`](tracks/systems/corpus/raw/gstelang__job-worker-service/) | 2024-08-26 → 2024-09-05 | unknown |  | 35 | **3** ⚠️ <sub>A:-1 B:7</sub> |
| [`zeeshanhaque21/JobWorkerService`](tracks/systems/corpus/raw/zeeshanhaque21__JobWorkerService/) | 2024-07-19 → 2024-08-03 | — |  | 37 | **3** <sub>A:4 B:2</sub> |
| [`bkneis/jobworker`](tracks/systems/corpus/raw/bkneis__jobworker/) | 2024-06-25 | — |  | 99 | **2.5** <sub>A:2 B:3</sub> |
| [`kiakeshmiri/process-runner`](tracks/systems/corpus/raw/kiakeshmiri__process-runner/) | 2024-09-18 → 2024-09-22 | — |  | 64 | **2.5** <sub>A:3 B:2</sub> |
| [`jkj/telechallenge`](tracks/systems/corpus/raw/jkj__telechallenge/) | — | — |  | — | **2** ⚙ <sub>A:2 B:2</sub> |
| [`njayp/jobber`](tracks/systems/corpus/raw/njayp__jobber/) | — | — |  | — | **2** ⚙ <sub>A:2 B:2</sub> |
| [`chintamanil/job-worker`](tracks/systems/corpus/raw/chintamanil__job-worker/) | 2026-02-20 | L4 |  | 62 | **1.5** <sub>A:1 B:2</sub> |
| [`kshi36/teleport-backend-challenge`](tracks/systems/corpus/raw/kshi36__teleport-backend-challenge/) | — | — |  | — | **1.5** ⚙ <sub>A:1 B:2</sub> |
| [`kurczynski/teleport-job-worker`](tracks/systems/corpus/raw/kurczynski__teleport-job-worker/) | 2024-05-22 → 2024-06-04 | L4 |  | 65 | **1.5** ⚠️ <sub>A:3 B:0</sub> |
| [`P-A-R-U-S/Go-Job-Worker-Service`](tracks/systems/corpus/raw/P-A-R-U-S__Go-Job-Worker-Service/) | 2024-08-25 → 2024-09-07 | — |  | 93 | **1.5** <sub>A:2 B:1</sub> |
| [`aadc-dev/teleport-job-management-int`](tracks/systems/corpus/raw/aadc-dev__teleport-job-management-int/) | — | — |  | — | **1** ⚙ <sub>A:1 B:1</sub> |
| [`bejelith/goteleport`](tracks/systems/corpus/raw/bejelith__goteleport/) | — | — |  | — | **1** ⚙ <sub>A:0 B:2</sub> |
| [`kovyrin/teleport-exec`](tracks/systems/corpus/raw/kovyrin__teleport-exec/) | 2022-01-29 → 2022-01-30 | — |  | 142 | **1** ⚠️ <sub>A:-1 B:3</sub> |
| [`wrboyce/jobsvc`](tracks/systems/corpus/raw/wrboyce__jobsvc/) | — | — |  | — | **1** ⚠️ ⚙ <sub>A:3 B:-1</sub> |
| [`bucknercd/jobworker`](tracks/systems/corpus/raw/bucknercd__jobworker/) | 2025-08-12 | L5 |  | 114 | **0** ⚠️ <sub>A:-2 B:2</sub> |
| [`kelwang/teleport`](tracks/systems/corpus/raw/kelwang__teleport/) | — | — |  | — | **0** ⚙ <sub>A:0 B:0</sub> |
| [`Kieran2k15/teleproc`](tracks/systems/corpus/raw/Kieran2k15__teleproc/) | — | — |  | — | **0** ⚙ <sub>A:0 B:0</sub> |
| [`nschneid95/TeleportInterview`](tracks/systems/corpus/raw/nschneid95__TeleportInterview/) | — | — |  | — | **0** ⚙ <sub>A:1 B:-1</sub> |
| [`xp10wa/job-worker-service`](tracks/systems/corpus/raw/xp10wa__job-worker-service/) | — | — |  | — | **0** ⚠️ ⚙ <sub>A:2 B:-2</sub> |
| [`derekschultz/job-worker`](tracks/systems/corpus/raw/derekschultz__job-worker/) | — | — |  | — | **-0.5** ⚠️ ⚙ <sub>A:-2 B:1</sub> |
| [`diptadas/job-worker`](tracks/systems/corpus/raw/diptadas__job-worker/) | 2021-06-05 → 2021-06-07 | L3 |  | 72 | **-0.5** ⚠️ <sub>A:-2 B:1</sub> |
| [`p4n1c/jws`](tracks/systems/corpus/raw/p4n1c__jws/) | — | — |  | — | **-0.5** ⚠️ ⚙ <sub>A:2 B:-3</sub> |
| [`parmjassal/assignment-job-worker`](tracks/systems/corpus/raw/parmjassal__assignment-job-worker/) | — | — |  | — | **-0.5** ⚙ <sub>A:0 B:-1</sub> |
| [`ChrisGe4/teleport`](tracks/systems/corpus/raw/ChrisGe4__teleport/) | — | — |  | — | **-1** ⚙ <sub>A:-2 B:0</sub> |
| [`cweiser22/teleport_interview`](tracks/systems/corpus/raw/cweiser22__teleport_interview/) | — | — |  | — | **-1** ⚙ <sub>A:-2 B:0</sub> |
| [`jaylane/job-scheduler`](tracks/systems/corpus/raw/jaylane__job-scheduler/) | 2024-08-09 → 2024-08-29 | unknown |  | 62 | **-1** ⚠️ <sub>A:-6 B:4</sub> |
| [`andrewhare/jobs`](tracks/systems/corpus/raw/andrewhare__jobs/) | — | — |  | — | **-3** ⚠️ ⚙ <sub>A:-9 B:3</sub> |
| [`jawnsy/sprinter`](tracks/systems/corpus/raw/jawnsy__sprinter/) | — | — |  | — | **-3** ⚠️ ⚙ <sub>A:0 B:-6</sub> |
| [`mcampo84/teleport_challenge`](tracks/systems/corpus/raw/mcampo84__teleport_challenge/) | 2025-01-08 → 2025-01-17 | — |  | 64 | **-3** ⚠️ <sub>A:-6 B:0</sub> |
| [`rafaelvanoni/jobservice`](tracks/systems/corpus/raw/rafaelvanoni__jobservice/) | — | — |  | — | **-3** ⚙ <sub>A:-3 B:-3</sub> |
| [`rexposadas/teleport`](tracks/systems/corpus/raw/rexposadas__teleport/) | 2024-10-07 → 2024-10-16 | — |  | 66 | **-3** <sub>A:-3 B:-3</sub> |
| [`xiaofatiandi/JobWorkerService`](tracks/systems/corpus/raw/xiaofatiandi__JobWorkerService/) | — | — |  | — | **-3** ⚠️ ⚙ <sub>A:-1 B:-5</sub> |
| [`samschurter/teleport-challenge`](tracks/systems/corpus/raw/samschurter__teleport-challenge/) | 2021-08-22 → 2021-08-31 | — |  | 59 | **-3.5** ⚠️ <sub>A:-8 B:1</sub> |
| [`devnulled/runjob`](tracks/systems/corpus/raw/devnulled__runjob/) | 2024-10-10 | — |  | 33 | **-4** ⚠️ <sub>A:-9 B:1</sub> |
| [`hashsequence/Linux-Job-Worker`](tracks/systems/corpus/raw/hashsequence__Linux-Job-Worker/) | 2020-08-26 | — |  | 65 | **-4** <sub>A:-4 B:-4</sub> |
| [`mikihau/rest-job-worker`](tracks/systems/corpus/raw/mikihau__rest-job-worker/) | — | unknown |  | — | **-4** ⚠️ <sub>A:-9 B:1</sub> |
| [`razzam21/job-worker-service`](tracks/systems/corpus/raw/razzam21__job-worker-service/) | 2025-08-07 | L4 |  | 65 | **-4** <sub>A:-3 B:-5</sub> |
| [`s-gruneberg/jobWorker`](tracks/systems/corpus/raw/s-gruneberg__jobWorker/) | 2025-08-27 → 2025-08-30 | — |  | 25 | **-4** ⚠️ <sub>A:-1 B:-7</sub> |
| [`clydotron/job_worker_service`](tracks/systems/corpus/raw/clydotron__job_worker_service/) | 2024-05-14 | — |  | 39 | **-4.5** ⚠️ <sub>A:-11 B:2</sub> |
| [`dakotasanchez/job-worker-challenge`](tracks/systems/corpus/raw/dakotasanchez__job-worker-challenge/) | — | — |  | — | **-4.5** ⚠️ ⚙ <sub>A:-9 B:0</sub> |
| [`moalf/teleport`](tracks/systems/corpus/raw/moalf__teleport/) | 2025-12-05 → 2025-12-17 | — |  | 73 | **-4.5** ⚠️ <sub>A:-13 B:4</sub> |
| [`lsiv568/jobs-worker`](tracks/systems/corpus/raw/lsiv568__jobs-worker/) | — | — |  | — | **-5** ⚠️ ⚙ <sub>A:-2 B:-8</sub> |
| [`zcancio/challenge1`](tracks/systems/corpus/raw/zcancio__challenge1/) | — | — |  | — | **-5** ⚠️ ⚙ <sub>A:-2 B:-8</sub> |
| [`bill-rich/jobworker`](tracks/systems/corpus/raw/bill-rich__jobworker/) | 2024-05-19 | — |  | 35 | **-5.5** ⚠️ <sub>A:-12 B:1</sub> |
| [`ehsaniara/joblet`](tracks/systems/corpus/raw/ehsaniara__joblet/) | 2025-06-19 → 2026-04-25 | — |  | 52 | **-7** ⚠️ <sub>A:-13 B:-1</sub> |
| [`m3talsmith/jobberthehut`](tracks/systems/corpus/raw/m3talsmith__jobberthehut/) | 2025-11-06 → 2025-11-07 | — |  | 91 | **-7** <sub>A:-7 B:-7</sub> |
| [`antonefremov/teleport_job_worker`](tracks/systems/corpus/raw/antonefremov__teleport_job_worker/) | — | — |  | — | **-7.5** ⚠️ <sub>A:-15 B:0</sub> |
| [`DiscoRiver/jobber`](tracks/systems/corpus/raw/DiscoRiver__jobber/) | 2025-03-11 | — |  | 65 | **-7.5** ⚠️ <sub>A:-12 B:-3</sub> |
| [`rsteinkeXJ/teleport_interview`](tracks/systems/corpus/raw/rsteinkeXJ__teleport_interview/) | 2026-04-22 | — |  | 77 | **-7.5** ⚠️ <sub>A:-10 B:-5</sub> |
| [`reynn/grpc-job-processor`](tracks/systems/corpus/raw/reynn__grpc-job-processor/) | — | — |  | — | **-8** ⚠️ ⚙ <sub>A:-13 B:-3</sub> |
| [`shawon-crosen/teleport-challenge`](tracks/systems/corpus/raw/shawon-crosen__teleport-challenge/) | 2025-06-03 | — |  | 51 | **-8** <sub>A:-9 B:-7</sub> |
| [`vincenthlam/teleport-interview`](tracks/systems/corpus/raw/vincenthlam__teleport-interview/) | 2024-05-23 | — |  | 7 | **-8** ⚠️ <sub>A:-15 B:-1</sub> |
| [`dtom90/teleport-backend-challenge`](tracks/systems/corpus/raw/dtom90__teleport-backend-challenge/) | 2025-04-21 | — |  | 50 | **-8.5** ⚠️ <sub>A:-15 B:-2</sub> |
| [`flychicken123/Job_Worker`](tracks/systems/corpus/raw/flychicken123__Job_Worker/) | 2025-04-29 | — |  | 58 | **-8.5** ⚠️ <sub>A:-15 B:-2</sub> |
| [`sabernabil12/teleport-job-worker`](tracks/systems/corpus/raw/sabernabil12__teleport-job-worker/) | 2025-06-24 | — |  | 54 | **-8.5** ⚠️ <sub>A:-10 B:-7</sub> |
| [`tembio/teleport`](tracks/systems/corpus/raw/tembio__teleport/) | — | — |  | — | **-8.5** ⚠️ ⚙ <sub>A:-11 B:-6</sub> |
| [`Takumi2008/go-run`](tracks/systems/corpus/raw/Takumi2008__go-run/) | — | — |  | — | **-10** ⚠️ ⚙ <sub>A:-12 B:-8</sub> |
| [`ghostsquad/prototype-job-worker`](tracks/systems/corpus/raw/ghostsquad__prototype-job-worker/) | 2024-03-19 | — |  | 27 | **-11** ⚠️ <sub>A:-14 B:-8</sub> |
| [`gaurav36/systemd-wrapper`](tracks/systems/corpus/raw/gaurav36__systemd-wrapper/) | — | — |  | — | **-13** ⚠️ <sub>A:-22 B:-4</sub> |
| [`espadolini/challenged`](tracks/systems/corpus/raw/espadolini__challenged__hired/) | — | — | ✓ [†](#hired-note) | — | — |
| [`lxea/teleport_challenge`](tracks/systems/corpus/raw/lxea__teleport_challenge__hired/) | — | — | ✓ [†](#hired-note) | — | — |
| [`rosstimothy/worker`](tracks/systems/corpus/raw/rosstimothy__worker__hired/) | — | — | ✓ [†](#hired-note) | — | — |
| [`tigrato/telep`](tracks/systems/corpus/raw/tigrato__telep__hired/) | — | — | ✓ [†](#hired-note) | — | — |

## Full-stack (4 submissions)

| Repo | Date range | Level | Hired | PR cmts | AI score |
|---|---|---|---|---|---|
| [`ibeckermayer/teleport-interview`](tracks/fullstack/corpus/raw/ibeckermayer__teleport-interview/) | 2021-01-05 → 2021-01-16 | — |  | 115 | **6.5** ⚠️ <sub>A:8 B:5</sub> |
| [`zship/teleport-challenge`](tracks/fullstack/corpus/raw/zship__teleport-challenge/) | 2021-09-28 → 2021-10-13 | — |  | 90 | **3.5** <sub>A:3 B:4</sub> |
| [`atburke/teleport_interview`](tracks/fullstack/corpus/raw/atburke__teleport_interview/) | 2021-05-31 | — |  | 22 | **3** <sub>A:4 B:2</sub> |
| [`tobocop/go-teleport-directory-browser`](tracks/fullstack/corpus/raw/tobocop__go-teleport-directory-browser/) | 2021-08-28 → 2021-08-31 | — |  | 55 | **2.5** <sub>A:3 B:2</sub> |

## Security automation (2 submissions)

| Repo | Date range | Level | Hired | PR cmts | AI score |
|---|---|---|---|---|---|
| [`tedmist1/challenge-auth0`](tracks/security-automation/corpus/raw/tedmist1__challenge-auth0/) | 2025-11-02 | — |  | 20 | **-1.5** ⚠️ <sub>A:-3 B:0</sub> |
| [`nick-hinds/teleport_challenge`](tracks/security-automation/corpus/raw/nick-hinds__teleport_challenge/) | 2025-09-19 | — |  | 26 | **-3.5** ⚠️ <sub>A:-7 B:0</sub> |

## SRE (1 submissions)

| Repo | Date range | Level | Hired | PR cmts | AI score |
|---|---|---|---|---|---|
| [`Chili-Man/teleport-sre-challenge`](tracks/sre/corpus/raw/Chili-Man__teleport-sre-challenge/) | 2023-12-03 → 2023-12-22 | — |  | 8 | **0.5** ⚠️ <sub>A:-1 B:2</sub> |

---

<a id="hired-note"></a>
**† Hired column**: A ✓ in this column means the submitter is a current Teleport reviewer in the captured PR review threads — i.e. they did the challenge themselves before joining and we have evidence they later reviewed others' submissions for the company. We do **not** know hire outcomes for the other submissions; absence of a ✓ is *not* a rejection signal. Hired-author rows show `—` for AI score because the source code is no longer accessible (repos were deleted or made private after the author joined Teleport); we have only the event metadata reconstructed from `data.gharchive.org`.

---

See also:

- [`SCORING_RUBRIC.md`](SCORING_RUBRIC.md) — how the AI score column was calculated
- [`mirrors/INDEX.md`](mirrors/INDEX.md) — read-only repo mirrors
- [`tracks/systems/`](tracks/systems/) — the most developed track, with reviewer-feedback corpus and pitfall docs
