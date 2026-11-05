const std = @import("std");

const Block = struct { mem: []u8, freed: bool };

fn nextSize(seed: *u32) usize {
    seed.* = seed.* *% 1103515245 +% 12345;
    return 4 + (seed.* >> 16) % 13;
}

pub fn main(init: std.process.Init) !void {
    const io = init.io;
    var out_buf: [2048]u8 = undefined;
    var in_buf: [256]u8 = undefined;
    var writer = std.Io.File.stdout().writerStreaming(io, &out_buf);
    var reader = std.Io.File.stdin().readerStreaming(io, &in_buf);
    const out = &writer.interface;
    const in = &reader.interface;

    var memory: [64]u8 = undefined;
    var fba = std.heap.FixedBufferAllocator.init(&memory);
    const alloc = fba.allocator();

    var blocks: [16]Block = undefined;
    var count: usize = 0;
    var served: usize = 0;
    var seed: u32 = 7;
    var want = nextSize(&seed);
    var msg_buf: [200]u8 = undefined;
    var msg: []const u8 = "";

    while (true) {
        var cells: [64]u8 = undefined;
        @memset(&cells, '.');
        var live: usize = 0;
        for (blocks[0..count], 0..) |b, i| {
            const off = @intFromPtr(b.mem.ptr) - @intFromPtr(&memory);
            if (!b.freed) live += b.mem.len;
            for (off..off + b.mem.len) |p| {
                if (b.freed) {
                    if (p < fba.end_index) cells[p] = 'x';
                } else cells[p] = 'A' + @as(u8, @intCast(i));
            }
        }

        try out.print("\x1b[H\x1b[2J", .{});

        for (0..8) |row| {
            for (cells[row * 8 ..][0..8]) |c| try out.print("{c} ", .{c});
            try out.print("\n", .{});
        }
        try out.print("\nЗанято живыми: {d} Б   bump-указатель: {d}/64 Б   запросов выполнено: {d}\n", .{ live, fba.end_index, served });

        if (msg.len == 0) {
            try out.print("\nВыделить {d} байт?\n", .{want});
        } else {
            try out.print("\nВыделить {d} байт?\n> {s}\n", .{ want, msg });
        }

        try out.print("\n[a] выделить  [f X] освободить блок X  [r] сброс арены  [q] выход\n> ", .{});
        try out.flush();

        const line = (try in.takeDelimiter('\n')) orelse break;
        const cmd = std.mem.trim(u8, line, " \r");
        if (cmd.len == 0) continue;
        switch (cmd[0]) {
            'a' => {
                if (count == blocks.len) {
                    msg = "таблица блоков заполнена";
                } else if (alloc.alloc(u8, want)) |buf| {
                    blocks[count] = .{ .mem = buf, .freed = false };
                    count += 1;
                    served += 1;
                    want = nextSize(&seed);
                    msg = "выделено успешно";
                } else |err| {
                    msg = std.fmt.bufPrint(&msg_buf, "error.{s}! (освобождение блока в середине не поможет: байты возвращают только последний блок или reset())", .{@errorName(err)}) catch "error";
                }
            },
            'f' => {
                const arg = std.mem.trim(u8, cmd[1..], " ");
                const id = if (arg.len > 0) std.ascii.toUpper(arg[0]) -% 'A' else 255;
                if (id >= count) {
                    msg = "нет такого блока";
                } else if (blocks[id].freed) {
                    msg = "ДВОЙНОЕ ОСВОБОЖДЕНИЕ! В режиме Debug аллокаторы Zig здесь вызвали бы panic.";
                } else {
                    alloc.free(blocks[id].mem);
                    blocks[id].freed = true;
                    msg = "освобождено";
                }
            },
            'r' => {
                fba.reset();
                count = 0;
                msg = "сброс! все блоки исчезли разом (старые указатели теперь «висячие»)";
            },
            'q' => break,
            else => msg = "неизвестная команда (используйте латинскую раскладку: a, f X, r, q)",
        }
    }

    var leaks: usize = 0;
    for (blocks[0..count]) |b| {
        if (!b.freed) leaks += 1;
    }

    try out.print("\nзапросов выполнено: {d}   утёкших блоков: {d}{s}\n", .{ served, leaks, if (leaks == 0) "  - чисто!" else "  - вы забыли их освободить!" });
    try out.flush();
}
