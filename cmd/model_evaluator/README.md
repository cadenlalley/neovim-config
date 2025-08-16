# Model Evaluator

## Evaluation Summary

The current model (gpt-4o-mini-2024-07-18) appears to be the best choice, offering:

- The highest accuracy scores
- Reasonable latency
- The lowest cost
- Consistent performance across different test cases

The newer models (like gpt-5-nano) show significantly higher latency and cost without providing better accuracy, while the gpt-4o-2024-11-20 model shows lower accuracy across the board.

### Accuracy:
* Highest Accuracy: `gpt-4o-mini-2024-07-18` (99.47% on apple-crumble)
* Most Consistent: `gpt-4o-mini-2024-07-18` (92.29% - 99.47% across tests)
* Lowest Accuracy: `gpt-4o-2024-11-20` (69.47% - 86.96% across tests)

### Latency:
* Fastest: `gpt-4o-2024-11-20` (5,488ms - 9,952ms)
* Slowest: `gpt-5-nano-2025-08-07` (50,115ms - 84,838ms)

### Cost Efficiency:
* Most Cost-Effective: `gpt-4o-mini-2024-07-18` ($0.0009 - $0.0013 per test)
* Most Expensive: `gpt-4o-2024-11-20` ($0.0173 - $0.0219 per test)

## Evaluation Results

### gpt-4o-mini-2024-07-18 (current)

| Source | Model | Prompt Tokens | Completion Tokens | Latency (ms) | Accuracy | Total Cost |
| --- | --- | --- | --- | --- | --- | --- |
| apple-crumble | `gpt-4o-mini-2024-07-18` | 3958 | 722 | 9905 | 0.994720 | 0.001027 |
| shrimp-tacos | `gpt-4o-mini-2024-07-18` | 5032 | 1033 | 13864 | 0.922946 | 0.001375 |
| french-onion-soup | `gpt-4o-mini-2024-07-18` | 3106 | 953 | 13392 | 0.917509 | 0.001038 |
| coconut-macaroons | `gpt-4o-mini-2024-07-18` | 3452 | 638 | 9359 | 0.963365 | 0.000901 |
| shake-shack-burgers | `gpt-4o-mini-2024-07-18` | 4038 | 1010 | 15831 | 0.924370 | 0.001212 |

### gpt-4o-2024-11-20

| Source | Model | Prompt Tokens | Completion Tokens | Latency (ms) | Accuracy | Total Cost |
| --- | --- | --- | --- | --- | --- | --- |
| apple-crumble | `gpt-4o-2024-11-20` | 3958 | 742 | 7173 | 0.869633 | 0.017315 |
| shrimp-tacos | `gpt-4o-2024-11-20` | 5032 | 932 | 9461 | 0.777242 | 0.021900 |
| french-onion-soup | `gpt-4o-2024-11-20` | 3106 | 1007 | 9952 | 0.724939 | 0.017835 |
| coconut-macaroons | `gpt-4o-2024-11-20` | 3452 | 581 | 5488 | 0.694701 | 0.014440 |
| shake-shack-burgers | `gpt-4o-2024-11-20` | 4038 | 822 | 8300 | 0.745231 | 0.018315 |

### gpt-4.1-mini-2025-04-14

| Source | Model | Prompt Tokens | Completion Tokens | Latency (ms) | Accuracy | Total Cost |
| --- | --- | --- | --- | --- | --- | --- |
| french-onion-soup | `gpt-4.1-mini-2025-04-14` | 3106 | 954 | 14697 | 0.939243 | 0.002769 |
| coconut-macaroons | `gpt-4.1-mini-2025-04-14` | 3452 | 633 | 9883 | 0.917208 | 0.002394 |
| shake-shack-burgers | `gpt-4.1-mini-2025-04-14` | 4038 | 1058 | 14619 | 0.928648 | 0.003308 |
| apple-crumble | `gpt-4.1-mini-2025-04-14` | 3958 | 788 | 12332 | 0.918641 | 0.002844 |
| shrimp-tacos | `gpt-4.1-mini-2025-04-14` | 5032 | 1080 | 16742 | 0.914210 | 0.003741 |

### gpt-5-nano-2025-08-07

| Source | Model | Prompt Tokens | Completion Tokens | Latency (ms) | Accuracy | Total Cost |
| --- | --- | --- | --- | --- | --- | --- |
| apple-crumble | `gpt-5-nano-2025-08-07` | 3954 | 9428 | 50115 | 0.872161 | 0.003969 |
| shrimp-tacos | `gpt-5-nano-2025-08-07` | 5028 | 17303 | 84838 | 0.791359 | 0.007173 |
| french-onion-soup | `gpt-5-nano-2025-08-07` | 3102 | 11539 | 53653 | 0.846010 | 0.004771 |
| coconut-macaroons | `gpt-5-nano-2025-08-07` | 3448 | 10469 | 50352 | 0.884788 | 0.004360 |
| shake-shack-burgers | `gpt-5-nano-2025-08-07` | 4034 | 14370 | 75040 | 0.832033 | 0.005950 |