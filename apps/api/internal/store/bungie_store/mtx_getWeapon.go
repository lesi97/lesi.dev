package bungie_store

import (
	"context"
	"fmt"
)

type weaponResult struct {
	weaponData *weaponData
	weaponPerks *filteredPerksResult
	err error
}

func (s *BungieStore) getWeapon(ctx context.Context, weaponHashID string, perkHashIDs []string) (*weaponResult, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan weaponResult, 2)

	go func() {
		weapon, err := s.getWeaponData(ctx, weaponHashID)
		if err != nil {
			fmt.Printf("ERROR: mtx_getWeapon: getWeaponData: %v\n", err)
			ch <- weaponResult{err: err}
			return
		}
		ch <- weaponResult{weaponData: weapon}
	}()


	go func() {
		perks, err := s.getWeaponPerks(ctx, perkHashIDs)
		if err != nil {
			fmt.Printf("ERROR: mtx_getWeapon: getWeaponPerks: %v\n", err)
			ch <- weaponResult{err: err}
			return
		}
		ch <- weaponResult{weaponPerks: perks}
	}()

	var final weaponResult

	for range 2 {
		res := <-ch
		if res.err != nil {
			cancel()
			return nil, res.err
		}
		if res.weaponData != nil {
			final.weaponData = res.weaponData
		}
		if res.weaponPerks != nil {
			final.weaponPerks = res.weaponPerks
		}
	}

	return &final, nil
	
}