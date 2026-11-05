use std::io::{self, BufRead, Write};

struct Treasure;

impl Drop for Treasure {
    fn drop(&mut self) {
        println!("  [Drop] Сокровище освобождено - владелец вышел из области видимости.");
    }
}

const NAMES: [char; 4] = ['A', 'B', 'C', 'D'];

const HELP: &str = "
------------------------------- КАК ИГРАТЬ -------------------------------------
Сокровище [T] лежит на станции A. Отполируйте его и доставьте на станцию D.

ТРИ ПРАВИЛА (их же проверяет компилятор Rust):
  1. ВЛАДЕЛЕЦ: у сокровища всегда ровно ОДИН владелец (станция с отметкой [T]).
  2. ЗАИМСТВОВАНИЕ: любая станция может взять сокровище у владельца
     во временное пользование, не становясь владельцем:
       &     общее заимствование = можно только смотреть. Таких может быть много.
       &mut  изменяемое заимствование = можно менять (полировать). Только ОДНО,
             и никаких других заимствований рядом.
  3. Пока активно ЛЮБОЕ заимствование, владелец не может передать или удалить
     сокровище.

КОМАНДЫ (X, Y - станции):
  move X Y    передать владение от X к Y
  borrow X Y  Y берёт у X как & (readonly)
  mut X Y     Y берёт у X как &mut (ссылка с возможностью записывания) и ПОЛИРУЕТ
  end Y       Y возвращает своё заимствование           (сделайте это до передачи)

ПОБЕДА: все три пункта чек-листа отмечены [x]. 3 ошибки = сборка не удалась.
--------------------------------------------------------------------------------";

struct Game {
    slots: [Option<Box<Treasure>>; 4],
    borrows: Vec<(usize, bool)>,
    polished: bool,
    errors: u32,
}

fn idx(s: Option<&&str>) -> Option<usize> {
    match s?.to_uppercase().chars().next()? {
        'A' => Some(0),
        'B' => Some(1),
        'C' => Some(2),
        'D' => Some(3),
        _ => None,
    }
}

fn check(b: bool) -> &'static str {
    if b {
        "x"
    } else {
        " "
    }
}

impl Game {
    fn owner(&self) -> Option<usize> {
        self.slots.iter().position(|s| s.is_some())
    }

    fn need_owner(&self, x: usize) -> Result<(), String> {
        if self.slots[x].is_some() {
            return Ok(());
        }
        Err(format!(
            "error[E0382]: использование перемещённого значения\n\t`{}` не владеет сокровищем. Одолжить или передать его может только владелец.",
            NAMES[x]
        ))
    }

    fn exec(&mut self, cmd: &[&str]) -> Result<String, String> {
        let (a, b) = (idx(cmd.get(1)), idx(cmd.get(2)));
        let is_borrowed = !self.borrows.is_empty();
        let moved_err = |x: usize| {
            format!("error[E0505]: нельзя переместить `{}`, пока оно заимствовано\n\tКто-то ещё держит заимствование. Сначала верните его командой `end <станция>`.", NAMES[x])
        };

        match (cmd[0], a, b) {
            ("move", Some(x), Some(y)) => {
                self.need_owner(x)?;
                if is_borrowed {
                    return Err(moved_err(x));
                }
                self.slots[y] = self.slots[x].take();
                Ok(format!("владение передано {} -> {} (станция {} теперь пуста)", NAMES[x], NAMES[y], NAMES[x]))
            }
            ("borrow", Some(x), Some(y)) => {
                self.need_owner(x)?;
                if self.borrows.iter().any(|b| b.1) {
                    return Err("error[E0502]: нельзя взять общее заимствование, пока есть изменяемое\n\tПока кто-то может менять значение (&mut), никто другой не может даже смотреть на него.".into());
                }
                self.borrows.push((y, false));
                Ok(format!("{} держит &Treasure (readonly). Владелец по-прежнему {}.", NAMES[y], NAMES[x]))
            }
            ("mut", Some(x), Some(y)) => {
                self.need_owner(x)?;
                if let Some(&(_, is_mut)) = self.borrows.first() {
                    return Err(if is_mut {
                        "error[E0499]: нельзя взять изменяемое заимствование дважды одновременно\n\t&mut может быть только ОДНО.".into()
                    } else {
                        "error[E0502]: нельзя взять изменяемое заимствование, пока есть общее\n\t&mut должно быть в одиночестве.".into()
                    });
                }
                self.borrows.push((y, true));
                self.polished = true;
                Ok(format!("{} держит &mut Treasure и ПОЛИРУЕТ его! Теперь верните: `end {}`", NAMES[y], NAMES[y]))
            }
            ("end", Some(y), _) => {
                let before = self.borrows.len();
                self.borrows.retain(|b| b.0 != y);
                if before == self.borrows.len() {
                    return Err(format!("У станции {} нет заимствований, завершать нечего.", NAMES[y]));
                }
                Ok(format!("{} вернула заимствование (время жизни ссылки закончилось)", NAMES[y]))
            }
            _ => Err("неизвестная или неполная команда. Введите `help`. (Станции A-D, латинская раскладка.)".into()),
        }
    }

    fn won(&self) -> bool {
        self.owner() == Some(3) && self.borrows.is_empty() && self.polished
    }

    fn show(&self) {
        let row: Vec<String> = (0..4)
            .map(|i| {
                let t = if self.slots[i].is_some() {
                    "[T]"
                } else {
                    "[ ]"
                };
                let b: Vec<&str> = self
                    .borrows
                    .iter()
                    .filter(|b| b.0 == i)
                    .map(|b| if b.1 { "&mut" } else { "&" })
                    .collect();
                format!("{}{} {}", NAMES[i], t, b.join(""))
            })
            .collect();

        println!("\n\t{}", row.join("    "));
        println!(
            "\tЦЕЛЬ: [{}] отполировано   [{}] нет активных заимствований   [{}] сокровище на D      ошибки: {}/3",
            check(self.polished), check(self.borrows.is_empty()), check(self.owner() == Some(3)), self.errors
        );
    }
}

fn main() {
    let mut g = Game {
        slots: [Some(Box::new(Treasure)), None, None, None],
        borrows: vec![],
        polished: false,
        errors: 0,
    };

    println!("{}", HELP);

    loop {
        g.show();
        if g.won() {
            println!("\nСборка успешна! Сокровище отполировано и доставлено.");
            break;
        }
        if g.owner().is_none() {
            println!("\nСокровище потеряно. Игра окончена.");
            break;
        }
        if g.errors >= 3 {
            println!(
                "\nerror: не удалось скомпилировать (3 ошибки borrow checker). Игра окончена."
            );
            break;
        }

        print!("> ");

        io::stdout().flush().unwrap();
        let mut line = String::new();
        if io::stdin().lock().read_line(&mut line).unwrap() == 0 {
            break;
        }

        let words: Vec<&str> = line.split_whitespace().collect();

        match words.first().copied() {
            None => continue,
            Some("quit") => break,
            _ => match g.exec(&words) {
                Ok(m) => println!("  ок: {}", m),
                Err(e) => {
                    g.errors += 1;
                    println!("  {}", e);
                }
            },
        }
    }

    println!("конец main(): переменная `g` выходит из области видимости...");
}
