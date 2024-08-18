package csw.chulbongkr.util;

import lombok.Getter;
import lombok.Setter;

import java.io.*;
import java.util.*;

/**
 * An implementation of Aho-Corasick algorithm based on Double Array Trie
 *
 * @param <V> the value type
 */
public class DAT<V> implements Serializable {
    @Serial
    private static final long serialVersionUID = 2L;

    /**
     * Check array of the Double Array Trie structure
     */
    private int[] check;
    /**
     * Base array of the Double Array Trie structure
     */
    private int[] base;
    /**
     * Fail table of the Aho-Corasick automaton
     */
    private int[] fail;
    /**
     * Output table of the Aho-Corasick automaton
     */
    private int[][] output;
    /**
     * Outer value array
     */
    private V[] values;
    /**
     * The length of each key
     */
    private int[] lengths;
    /**
     * The size of base and check arrays
     */
    private int size;

    /**
     * Parse text
     *
     * @param text The text
     * @return a list of outputs
     */
    public List<Hit<V>> parseText(CharSequence text) {
        return parseTextInternal(text, (begin, end, value) -> true);
    }

    /**
     * Parse text with a processor
     *
     * @param text      The text
     * @param processor A processor which handles the output
     */
    public void processText(CharSequence text, IHit<V> processor) {
        parseTextInternal(text, (begin, end, value) -> {
            processor.hit(begin, end, value);
            return true;
        });
    }

    /**
     * Parse text with a cancellable processor
     *
     * @param text      The text
     * @param processor A processor which handles the output
     */
    public void processTextCancellable(CharSequence text, IHitCancellable<V> processor) {
        parseTextInternal(text, processor);
    }

    /**
     * Common parsing method with a cancellable function
     *
     * @param text      The text
     * @param processor A processor function which handles the output
     */
    private List<Hit<V>> parseTextInternal(CharSequence text, IHitCancellable<V> processor) {
        List<Hit<V>> collectedEmits = new ArrayList<>();
        int currentState = 0;
        for (int position = 0; position < text.length(); ++position) {
            currentState = getState(currentState, text.charAt(position));
            int[] hitArray = output[currentState];
            if (hitArray != null) {
                for (int hit : hitArray) {
                    int begin = position + 1 - lengths[hit];
                    int end = position + 1;
                    V value = values[hit];
                    if (!processor.hit(begin, end, value)) {
                        return collectedEmits;
                    }
                    collectedEmits.add(new Hit<>(begin, end, value));
                }
            }
        }
        return collectedEmits;
    }

    /**
     * Checks if the text contains at least one substring
     *
     * @param text source text to check
     * @return {@code true} if the string contains at least one substring
     */
    public boolean matches(CharSequence text) {
        return findFirst(text) != null;
    }

    /**
     * Search for the first match in the text
     *
     * @param text source text to check
     * @return first match or {@code null} if there are no matches
     */
    public Hit<V> findFirst(CharSequence text) {
        int currentState = 0;
        int textLength = text.length();

        for (int position = 0; position < textLength; ++position) {
            currentState = getState(currentState, text.charAt(position));
            int[] hitArray = output[currentState];
            if (hitArray != null) {
                int hitIndex = hitArray[0];
                return new Hit<>(position + 1 - lengths[hitIndex], position + 1, values[hitIndex]);
            }
        }
        return null;
    }

    /**
     * Save the trie to an ObjectOutputStream
     *
     * @param out An ObjectOutputStream object
     * @throws IOException Some IOException
     */
    public void save(ObjectOutputStream out) throws IOException {
        out.writeObject(base);
        out.writeObject(check);
        out.writeObject(fail);
        out.writeObject(output);
        out.writeObject(lengths);
        out.writeObject(values);
    }

    /**
     * Load data from an ObjectInputStream
     *
     * @param in An ObjectInputStream object
     * @throws IOException            If can't read the file from path
     * @throws ClassNotFoundException If the class doesn't exist or match
     */
    @SuppressWarnings("unchecked")
    public void load(ObjectInputStream in) throws IOException, ClassNotFoundException {
        base = (int[]) in.readObject();
        check = (int[]) in.readObject();
        fail = (int[]) in.readObject();
        output = (int[][]) in.readObject();
        lengths = (int[]) in.readObject();
        values = (V[]) in.readObject();
    }

    /**
     * Get value by a String key, similar to a map.get() method
     *
     * @param key The key
     * @return value if it exists, otherwise returns null
     */
    public V get(CharSequence key) {
        int index = exactMatchSearch(key);
        return index >= 0 ? values[index] : null;
    }

    /**
     * Update a value corresponding to a key
     *
     * @param key   the key
     * @param value the value
     * @return successful or not (failure if there is no key)
     */
    public boolean set(CharSequence key, V value) {
        int index = exactMatchSearch(key);
        if (index >= 0) {
            values[index] = value;
            return true;
        }
        return false;
    }

    /**
     * Pick the value by index in value array
     *
     * @param index The index
     * @return The value
     */
    public V get(int index) {
        return values[index];
    }

    /**
     * Processor handles the output when hitting a keyword
     */
    public interface IHit<V> {
        void hit(int begin, int end, V value);
    }

    /**
     * Processor handles the output when hitting a keyword, with more detail
     */
    public interface IHitFull<V> {
        void hit(int begin, int end, V value, int index);
    }

    /**
     * Callback that allows cancelling the search process
     */
    public interface IHitCancellable<V> {
        boolean hit(int begin, int end, V value);
    }

    /**
         * A result output
         *
         * @param <V> the value type
         */
        public record Hit<V>(int begin, int end, V value) {

        @Override
            public String toString() {
                return String.format("[%d:%d]=%s", begin, end, value);
            }
        }

    /**
     * Transmit state, supports failure function
     *
     * @param currentState The current state
     * @param character    The character
     * @return The new state
     */
    private int getState(int currentState, char character) {
        int newState = transitionWithRoot(currentState, character);
        while (newState == -1) {
            currentState = fail[currentState];
            newState = transitionWithRoot(currentState, character);
        }
        return newState;
    }

    /**
     * Transition of a state
     *
     * @param current The current state
     * @param c       The character
     * @return The new state
     */
    protected int transition(int current, char c) {
        int p = current + c + 1;
        return current == check[p] ? base[p] : -1;
    }

    /**
     * Transition of a state, if the state is root and it failed, then returns the root
     *
     * @param nodePos The node position
     * @param c       The character
     * @return The new state
     */
    protected int transitionWithRoot(int nodePos, char c) {
        int b = base[nodePos];
        int p = b + c + 1;
        return b != check[p] ? (nodePos == 0 ? 0 : -1) : p;
    }

    /**
     * Build an AhoCorasickDoubleArrayTrie from a map
     *
     * @param map A map containing key-value pairs
     */
    public void build(Map<String, V> map) {
        new Builder().build(map);
    }

    /**
     * Match exactly by a key
     *
     * @param key The key
     * @return The index of the key, which can be used as a perfect hash function
     */
    public int exactMatchSearch(CharSequence key) {
        return exactMatchSearch(key, 0, key.length(), 0);
    }

    /**
     * Match exactly by a key
     *
     * @param key      The key
     * @param pos      The position
     * @param len      The length
     * @param nodePos  The node position
     * @return The index of the key
     */
    private int exactMatchSearch(CharSequence key, int pos, int len, int nodePos) {
        int b = base[nodePos];
        int p;

        for (int i = pos; i < len; i++) {
            p = b + key.charAt(i) + 1;
            if (b != check[p]) {
                return -1;
            }
            b = base[p];
        }

        p = b;
        int n = base[p];
        return b == check[p] ? -n - 1 : -1;
    }

    /**
     * @return the size of the keywords
     */
    public int size() {
        return values.length;
    }

    /**
     * A builder to build the AhoCorasickDoubleArrayTrie
     */
    private class Builder {
        private final State rootState = new State();
        private BitSet used;
        private int allocSize;
        private int progress;
        private int nextCheckPos;
        private int keySize;

        public void build(Map<String, V> map) {
            if (map.isEmpty()) {
                throw new IllegalArgumentException("The input map is empty.");
            }

            values = (V[]) map.values().toArray();
            lengths = new int[values.length];
            Set<String> keySet = map.keySet();
            addAllKeywords(keySet);
            buildDoubleArrayTrie(keySet.size());
            used = null;
            constructFailureStates();
            reduceSize();
        }

        private void addAllKeywords(Collection<String> keywordSet) {
            int i = 0;
            for (String keyword : keywordSet) {
                addKeyword(keyword, i++);
            }
        }

        private void addKeyword(String keyword, int index) {
            State currentState = this.rootState;
            for (char character : keyword.toCharArray()) {
                currentState = currentState.addState(character);
            }
            currentState.addEmit(index);
            lengths[index] = keyword.length();
        }

        private void buildDoubleArrayTrie(int keySize) {
            progress = 0;
            this.keySize = keySize;
            resize(65536 * 32);

            base[0] = 1;
            nextCheckPos = 0;

            var siblings = new ArrayList<Map.Entry<Integer, State>>(rootState.getSuccess().size());
            fetch(rootState, siblings);
            if (!siblings.isEmpty()) {
                insert(siblings);
            }
        }

        private void resize(int newSize) {
            int[] newBase = new int[newSize];
            int[] newCheck = new int[newSize];

            if (allocSize > 0) {
                System.arraycopy(base, 0, newBase, 0, allocSize);
                System.arraycopy(check, 0, newCheck, 0, allocSize);
            }
            base = newBase;
            check = newCheck;
            used = new BitSet(newSize);
            allocSize = newSize;
        }

        private void insert(List<Map.Entry<Integer, State>> siblings) {
            var queue = new ArrayDeque<Map.Entry<Integer, List<Map.Entry<Integer, State>>>>();
            queue.add(new AbstractMap.SimpleEntry<>(null, siblings));

            while (!queue.isEmpty()) {
                var entry = queue.remove();
                var currentSiblings = entry.getValue();

                int begin = 0;
                int pos = Math.max(currentSiblings.getFirst().getKey() + 1, nextCheckPos) - 1;
                int nonzeroCount = 0;
                int first = 0;

                if (allocSize <= pos) {
                    resize(pos + 1);
                }

                while (true) {
                    pos++;
                    if (allocSize <= pos) {
                        resize(pos + 1);
                    }

                    if (check[pos] != 0) {
                        nonzeroCount++;
                        continue;
                    } else if (first == 0) {
                        nextCheckPos = pos;
                        first = 1;
                    }

                    begin = pos - currentSiblings.getFirst().getKey();
                    if (allocSize <= (begin + currentSiblings.getLast().getKey())) {
                        double toSize = Math.max(1.05, 1.0 * keySize / (progress + 1)) * allocSize;
                        int maxSize = (int) (Integer.MAX_VALUE * 0.95);
                        if (allocSize >= maxSize) throw new RuntimeException("Double array trie is too big.");
                        resize((int) Math.min(toSize, maxSize));
                    }

                    if (used.get(begin)) continue;

                    int finalBegin1 = begin;
                    if (currentSiblings.stream().anyMatch(sibling -> check[finalBegin1 + sibling.getKey()] != 0)) {
                        continue;
                    }

                    break;
                }

                if (1.0 * nonzeroCount / (pos - nextCheckPos + 1) >= 0.95) {
                    nextCheckPos = pos;
                }
                used.set(begin);

                size = Math.max(size, begin + currentSiblings.getLast().getKey() + 1);

                int finalBegin = begin;
                currentSiblings.forEach(sibling -> check[finalBegin + sibling.getKey()] = finalBegin);

                for (var sibling : currentSiblings) {
                    var newSiblings = new ArrayList<Map.Entry<Integer, State>>(sibling.getValue().getSuccess().size() + 1);

                    if (fetch(sibling.getValue(), newSiblings) == 0) {
                        base[begin + sibling.getKey()] = -sibling.getValue().getLargestValueId() - 1;
                        progress++;
                    } else {
                        queue.add(new AbstractMap.SimpleEntry<>(begin + sibling.getKey(), newSiblings));
                    }
                    sibling.getValue().setIndex(begin + sibling.getKey());
                }

                Integer parentBaseIndex = entry.getKey();
                if (parentBaseIndex != null) {
                    base[parentBaseIndex] = begin;
                }
            }
        }


        private void constructFailureStates() {
            fail = new int[size + 1];
            output = new int[size + 1][];
            var queue = new ArrayDeque<State>();

            rootState.getStates().forEach(depthOneState -> {
                depthOneState.setFailure(rootState, fail);
                queue.add(depthOneState);
                constructOutput(depthOneState);
            });

            while (!queue.isEmpty()) {
                State currentState = queue.remove();

                for (Character transition : currentState.getTransitions()) {
                    State targetState = currentState.nextState(transition);
                    queue.add(targetState);

                    State traceFailureState = currentState.failure();
                    while (traceFailureState != null && traceFailureState.nextState(transition) == null) {
                        traceFailureState = traceFailureState.failure();
                    }
                    State newFailureState = traceFailureState == null ? rootState : traceFailureState.nextState(transition);
                    targetState.setFailure(newFailureState, fail);
                    targetState.addEmit(newFailureState.emit());
                    constructOutput(targetState);
                }
            }
        }

        private void constructOutput(State targetState) {
            Collection<Integer> emit = targetState.emit();
            if (emit == null || emit.isEmpty()) return;
            int[] output = new int[emit.size()];
            int i = 0;
            for (int value : emit) {
                output[i++] = value;
            }
            DAT.this.output[targetState.getIndex()] = output;
        }

        private void reduceSize() {
            int newSize = size + 65535;
            base = Arrays.copyOf(base, newSize);
            check = Arrays.copyOf(check, newSize);
        }

        private int fetch(State parent, List<Map.Entry<Integer, State>> siblings) {
            if (parent.isAcceptable()) {
                State fakeNode = new State(-(parent.getDepth() + 1));
                fakeNode.addEmit(parent.getLargestValueId());
                siblings.add(new AbstractMap.SimpleEntry<>(0, fakeNode));
            }
            parent.getSuccess().forEach((key, value) -> siblings.add(new AbstractMap.SimpleEntry<>(key + 1, value)));
            return siblings.size();
        }
    }

    private static class State {
        @Getter
        private final int depth;

        @Getter
        private final Map<Character, State> success = new HashMap<>();
        private Set<Integer> emits = null;
        private State failure = null;
        @Getter
        private int largestValueId = -1;
        @Getter
        @Setter
        private int index;

        public State() {
            this(0);
        }

        public State(int depth) {
            this.depth = depth;
        }

        public State addState(Character character) {
            return success.computeIfAbsent(character, c -> new State(depth + 1));
        }

        public State nextState(Character character) {
            return success.get(character);
        }

        public Collection<State> getStates() {
            return success.values();
        }

        public Collection<Character> getTransitions() {
            return success.keySet();
        }

        public void addEmit(int keyword) {
            if (emits == null) emits = new HashSet<>();
            emits.add(keyword);
            if (keyword > largestValueId) largestValueId = keyword;
        }

        public void addEmit(Collection<Integer> emits) {
            if (this.emits == null) this.emits = new HashSet<>();
            this.emits.addAll(emits);
        }

        public Collection<Integer> emit() {
            return emits == null ? Collections.emptyList() : emits;
        }

        public boolean isAcceptable() {
            return depth > 0 && emits != null;
        }

        public State failure() {
            return failure;
        }

        public void setFailure(State failState, int[] fail) {
            this.failure = failState;
            fail[index] = failState.getIndex();
        }

    }
}
