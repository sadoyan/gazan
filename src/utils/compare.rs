use crate::utils::tools::*;
use tokio::sync::RwLockReadGuard;

/*
#[allow(dead_code)]
pub fn dashmaps(map1: &RwLockWriteGuard<UpstreamMap>, map2: &UpstreamMap) -> bool {
    if map1.len() != map2.len() {
        return false;
    }
    for entry1 in map1.iter() {
        let key = entry1.key();
        let (vec1, _) = entry1.value();

        if let Some(entry2) = map2.get(key) {
            let (vec2, _) = entry2.value();
            if vec1 != vec2 {
                return false;
            }
        } else {
            return false;
        }
    }
    true
}
*/

#[allow(dead_code)]
pub fn dm(map1: &RwLockReadGuard<UpstreamMap>, map2: &UpstreamMap) -> bool {
    if map1.len() != map2.len() {
        return false; // Different number of keys
    }
    for entry1 in map1.iter() {
        let key = entry1.key();
        let (vec1, _) = entry1.value(); // Extract value

        if let Some(entry2) = map2.get(key) {
            let (vec2, _) = entry2.value(); // Correctly extract value
            if vec1 != vec2 {
                return false;
            }
        } else {
            return false;
        }
    }
    true
}
