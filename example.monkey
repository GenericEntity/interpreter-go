/**
 * makeHash returns a hashmap from an input array of keys and an input array of values.
 * Both arrays should be of the same length.
 */
let makeHash = fn(keys, values) {
    if (len(keys) != len(values)) {
        puts("len(keys) != len(values)")
        return {}
    }

    if (len(keys) == 0) {
        return {}
    }

    return put(makeHash(rest(keys), rest(values)), first(keys), first(values))
};

puts(">> makeHash([1, 2, 3], [\"one\", \"two\", \"three\"])")
puts(makeHash([1, 2, 3], ["one", "two", "three"]));
puts("");

/**
 * map takes a mapping function, f, and applies it to each element in the array, xs,
 * then returns the resulting array.
 */
let map = fn(f, xs) {
    let map_iter = fn(f, xs, acc) {
        if (len(xs) == 0) {
            return acc
        }

        return map_iter(f, rest(xs), push(acc, f(first(xs))))
    }

    return map_iter(f, xs, [])
}

puts(">> map(fn(x){x*x}, [1,2,3])")
puts(map(fn(x){x*x}, [1,2,3]))
puts("");


/**
 * filter takes a predicate, keep, and tests each element in the array, xs,
 * returning an array containing each element in xs for which keep returns true.
 */
let filter = fn(keep, xs) {
    let filter_iter = fn(keep, xs, acc) {
        if (len(xs) == 0) {
            return acc
        }
        let curr = first(xs)
        let remaining = rest(xs)

        if (keep(curr)) {
            return filter_iter(keep, remaining, push(acc, curr))
        } else {
            return filter_iter(keep, remaining, acc)
        }
    }

    return filter_iter(keep, xs, [])
}

puts(">> filter(fn(x){ x > 10 }, [1,5,100, 29, 321])")
puts(filter(fn(x){ x > 10 }, [1,5,100, 29, 321]))
puts("");

/**
 * fold_right takes a binary operator, f; an array, xs; and an initial value, then
 * repeatedly applies f to combine the elements of xs from the right, then returns the result.
 * e.g. fold_right(fn(x, y){x + y}, [1,2,3,4], 0) performs (1 + (2 + (3 + (4 + 0))))
 */
let fold_right = fn(f, xs, initial) {
    if (len(xs) == 0) {
        return initial
    }

    let sp = fold_right(f, rest(xs), initial)

    return f(sp, first(xs))
}

puts(">> fold_right(fn(x, y) {x + y}, [1,2,3,4], 0)")
puts(
    fold_right(
        fn(x, y) {x + y},
        [1,2,3,4],
        0
    )
)
puts("");
