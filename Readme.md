
# AI Workflow – CV Recruiter System

## Overview

This system automatically ranks job applicants using AI.
When users apply for a job, their CVs are processed, analyzed, and compared against the job requirements.
The system then returns the top 5 best-matching candidates.

---

## High-Level Flow

1. User applies to a job
2. CV is processed and structured
3. Job description is processed
4. AI compares CVs with the job
5. Candidates are scored
6. Top 5 candidates are returned

---

## Step 1: Application Submission

When a user applies:

* Application is created
* CV is retrieved
* Application is queued for AI processing (background worker)

This ensures the API remains fast and non-blocking.

---

## Step 2: CV Processing

Each CV goes through:

1. Text Extraction

   * Extract raw text from uploaded file (PDF/DOCX)

2. Text Cleaning

   * Remove noise
   * Normalize formatting
   * Lowercase conversion

3. Information Extraction

   * Skills
   * Experience
   * Education
   * Certifications (optional)

4. Semantic Embedding Generation

   * Convert CV content into vector representation
   * Stored for similarity comparison

5. Structured Data Storage

   * Save parsed data + embedding for later matching

---

## Step 3: Job Processing

When a job is created or updated:

1. Extract required skills
2. Extract experience level
3. Extract education requirements
4. Generate semantic embedding from job description
5. Store structured job data + embedding

---

## Step 4: Candidate Matching

When requesting top candidates for a job:

For each applicant:

1. Retrieve processed CV data

2. Retrieve processed job data

3. Compute:

   * Skill Match Score
   * Experience Match Score
   * Education Match Score
   * Semantic Similarity Score (vector similarity)

4. Combine scores into a final weighted score

Example logic:

Final Score =
(0.4 × Skill Score) +
(0.2 × Experience Score) +
(0.1 × Education Score) +
(0.3 × Semantic Similarity)

All scores are stored for transparency and auditing.

---

## Step 5: Ranking

1. Sort candidates by final score
2. Select top 5
3. Return ranked list

Optional:

* Re-rank top candidates using LLM refinement
* Generate match explanation for recruiters

---

## Architecture Notes

* CV processing runs in background workers
* AI logic is isolated from core business logic
* Embeddings are stored for fast vector similarity comparison
* Scoring is deterministic and explainable

---

## Result

The system provides:

* Automated candidate ranking
* Objective scoring
* Scalable matching
* Fast recruiter decision support

