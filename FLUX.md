# Transformer Flow — Example Using the Sentence `"Java is good"`

This document explains, in a clear and detailed way, how a Transformer processes the sequence `"Java is good"` during the **forward pass**, computes the **loss**, performs the **backward pass**, and updates its weights.

---

## Input

Input sentence:

```text
"Java is good"
```

---

## 1. Tokenization

The sentence is split into tokens, and each token is converted into a numerical identifier (`token_id`):

* `"Java"` → `token_id(42)`
* `"is"` → `token_id(15)`
* `"good"` → `token_id(87)`

Tokenized sequence:

```text
[42, 15, 87]
```

---

## 2. Embedding

Each token is converted into a dense vector of size 64.

Since the sequence contains 3 tokens, the input matrix `X` has shape:

```text
X = [3, 64]
```

Where:

* `3` = number of tokens in the sentence
* `64` = embedding dimension

---

## 3. Positional Encoding

Because Transformers do not process tokens sequentially like RNNs, positional information must be added to the embeddings.

This is done by summing the embedding matrix with a positional encoding matrix:

```text
X = X + positional_encoding
```

The shape remains:

```text
[3, 64]
```

---

## 4. Forward Pass

The input passes through a stack of Transformer blocks.

In this example, there are **2 blocks**.

---

## Block 1 (of 2)

### 4.1 Multi-Head Attention

First, the model computes the **Query (Q)**, **Key (K)**, and **Value (V)** matrices:

```text
Q = X @ WQ → [3, 64]
K = X @ WK → [3, 64]
V = X @ WV → [3, 64]
```

Where:

* `@` = matrix multiplication
* `WQ`, `WK`, `WV` = learnable weight matrices

Next, attention scores are computed:

```text
scores = Q @ K.T → [3, 3]
```

This gives the similarity between each token and every other token.

The scores are then scaled:

```text
scores = scores / sqrt(64) → [3, 3]
```

After scaling, the softmax function converts the scores into attention weights:

```text
weights = softmax(scores) → [3, 3]
```

These weights are used to combine the value vectors:

```text
attention = weights @ V → [3, 64]
```

Then the attention output is projected again:

```text
attention = attention @ WO → [3, 64]
```

Where `WO` is another learnable matrix.

Finally, a residual connection and layer normalization are applied:

```text
X = layerNorm(X + attention) → [3, 64]
```

---

### 4.2 Feed-Forward Network (FFN)

After attention, the output goes through a feed-forward neural network.

First projection:

```text
H = X @ W1 → [3, 256]
```

Activation function:

```text
H = ReLU(H) → [3, 256]
```

Second projection:

```text
ffn = H @ W2 → [3, 64]
```

Then another residual connection and layer normalization:

```text
X = layerNorm(X + ffn) → [3, 64]
```

---

## Block 2 (of 2)

The same sequence of operations is repeated as in Block 1, but using a different set of learned weights.

At the end of Block 2, the output `X` still has shape:

```text
[3, 64]
```

---

## 5. Output Layer

The final hidden representation is projected into the vocabulary space.

```text
logits = X @ W_out → [3, 5000]
```

Where:

* `5000` = vocabulary size
* each row contains the raw prediction scores for the next token

---

## 6. Loss Computation

Assume the target sequence is:

```text
targets = [15, 87, 2]
```

Where:

* after `"Java"`, the expected token is `"is"` (`15`)
* after `"is"`, the expected token is `"good"` (`87`)
* after `"good"`, the expected token is the end token (`2`)

The loss is computed using cross-entropy:

```text
loss = cross_entropy(logits, targets)
```

This measures how far the model’s predictions are from the expected targets.

---

## 7. Backward Pass

The backward pass computes gradients for all parameters so the model can learn.

---

### 7.1 Loss Gradient

First, compute the gradient of the loss with respect to the logits:

```text
grad_logits = probs - one_hot(target) → [3, 5000]
```

---

### 7.2 Output Layer Backward

Gradient with respect to the output weights:

```text
grad_W_out = X.T @ grad_logits → [64, 5000]
```

Gradient passed back into the Transformer output:

```text
grad_X = grad_logits @ W_out.T → [3, 64]
```

---

### 7.3 Feed-Forward Network Backward (per block)

Gradient through the second FFN projection:

```text
grad_H = grad_X @ W2.T → [3, 256]
```

Gradient through the ReLU activation:

```text
grad_H = grad_H * (H > 0) → [3, 256]
```

Gradient for `W2`:

```text
grad_W2 = H.T @ grad_X → [256, 64]
```

Gradient for `W1`:

```text
grad_W1 = X.T @ grad_H → [64, 256]
```

Gradient propagated back to the block input:

```text
grad_X = grad_H @ W1.T → [3, 64]
```

---

### 7.4 Attention Backward (per block)

Gradient enters the attention output:

```text
grad_attention = grad_X → [3, 64]
```

Gradient for the output projection `WO`:

```text
grad_WO = attention_input.T @ grad_attention → [64, 64]
```

Gradient with respect to `V`:

```text
grad_V = weights.T @ grad_attention → [3, 64]
```

Gradient for `WV`:

```text
grad_WV = X.T @ grad_V → [64, 64]
```

Gradient with respect to the attention weights:

```text
grad_weights = grad_attention @ V.T → [3, 3]
```

Gradient through the softmax:

```text
grad_scores = grad_weights * softmax_derivative → [3, 3]
```

Gradient with respect to `Q`:

```text
grad_Q = grad_scores @ K → [3, 64]
```

Gradient for `WQ`:

```text
grad_WQ = X.T @ grad_Q → [64, 64]
```

Gradient with respect to `K`:

```text
grad_K = grad_scores.T @ Q → [3, 64]
```

Gradient for `WK`:

```text
grad_WK = X.T @ grad_K → [64, 64]
```

Combined gradient flowing back into the block input:

```text
grad_X = (grad_Q @ WQ.T) + (grad_K @ WK.T) + (grad_V @ WV.T) → [3, 64]
```

---

### 7.5 Embedding Backward

The gradient reaches the embedding layer:

```text
grad_embedding = grad_X → [3, 64]
```

The embedding vectors for the corresponding token IDs are updated:

```text
embedding_weights[token_id] -= lr * grad_embedding
```

Where `lr` is the learning rate.

---

## 8. Weight Update

After all gradients are computed, the model updates its parameters:

```text
W = W - lr * grad_W
```

This is done for all learnable matrices, including:

* `WQ`
* `WK`
* `WV`
* `WO`
* `W1`
* `W2`
* `W_out`
* embedding weights

---

## 9. Next Epoch

The process repeats with the next batch of training data:

1. tokenize input
2. compute embeddings
3. apply positional encoding
4. run the forward pass
5. compute loss
6. run the backward pass
7. update weights

Over many iterations, the model gradually improves its predictions.

---

## Summary

This example shows the complete training flow of a Transformer:

* **Tokenization** converts words into token IDs
* **Embeddings** transform tokens into vectors
* **Positional encoding** adds order information
* **Attention** lets tokens interact with one another
* **Feed-forward layers** refine representations
* **Output projection** maps hidden states to vocabulary scores
* **Cross-entropy loss** measures prediction error
* **Backpropagation** computes gradients
* **Weight updates** improve the model over time

This cycle is repeated across many batches and epochs until the model learns meaningful language patterns.
