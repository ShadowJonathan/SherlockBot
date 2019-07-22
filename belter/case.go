package Belt

import ()

func HandleMention(Men *FullMention, biChange *CompiledChange) ([]*PrimeSuspectChange, bool) {
	PS := sh.PrimeSuspects
	var ALLPSC []*PrimeSuspectChange
	Gold := biChange.Old
	Gnew := biChange.New
	var is bool
	for _, Ps := range PS {
		var PSC = &PrimeSuspectChange{
			ID: Ps,
		}
		PsMember := GetMember(Ps, Gnew)
		for _, C := range Men.ChannelOR {
			if !C.ExistCrisis {
				Cold := GetChannel(C.ID, Gold)
				Cnew := GetChannel(C.ID, Gnew)
				if C.perms {
					for _, Or := range C.Perms {
						if !Or.ExistCrisis {
							OOR := GetOR(Or.ID, Cold)
							NOR := GetOR(Or.ID, Cnew)
							if OOR.ID == Ps || NOR.ID == Ps || HasRole(PsMember, Or.ID) {
								var c bool
								var MEEM = &PSchannelchange{}
								if Or.Allow {
									MEEM.ID = OOR.ID
									MEEM.allow = true
									MEEM.Allow.Last = OOR.Allow
									MEEM.Allow.New = NOR.Allow
									PSC.channels = true
									c = true
								}
								if Or.Deny {
									MEEM.ID = OOR.ID
									MEEM.deny = true
									MEEM.Deny.Last = OOR.Deny
									MEEM.Deny.New = NOR.Deny
									PSC.channels = true
									c = true
								}
								if c {
									AllCC := PSC.CC
									AllCC = append(AllCC, MEEM)
									PSC.CC = AllCC
								}
							}
						} else {
							var MEEM = &PSchannelchange{}
							if Or.Mk {
								MEEM.ID = Or.ID
								MEEM.Made = true
								PSC.channels = true
								MEEM.Mper.Allowed = GetOR(Or.ID, Cnew).Allow
								MEEM.Mper.Denied = GetOR(Or.ID, Cnew).Deny
								MEEM.Existentcrisis = true
								MEEM.Type = GetOR(Or.ID, Cnew).Type

								AllCC := PSC.CC
								AllCC = append(AllCC, MEEM)
								PSC.CC = AllCC
							} else if Or.Del {
								MEEM.ID = Or.ID
								MEEM.Deleted = true
								PSC.channels = true
								MEEM.Dperm.Allowed = GetOR(Or.ID, Cold).Allow
								MEEM.Dperm.Denied = GetOR(Or.ID, Cold).Deny
								MEEM.Existentcrisis = true
								MEEM.Type = GetOR(Or.ID, Cold).Type

								AllCC := PSC.CC
								AllCC = append(AllCC, MEEM)
								PSC.CC = AllCC
							}
						}
					}
				}
			}
		}
		for _, R := range Men.Roles {
			if !R.ExistCrisis && HasRole(PsMember, R.ID) {
				var c bool
				var ROOR = &PSRolechange{}
				Ro := GetRole(R.ID, Gold)
				Rn := GetRole(R.ID, Gnew)
				if R.Perms {
					ROOR.ID = R.ID
					ROOR.perms = true
					ROOR.oldPerms = Ro.Permissions
					ROOR.newPerms = Rn.Permissions
					PSC.roles = true
					c = true
				}
				if R.Position {
					ROOR.ID = R.ID
					ROOR.pos = true
					ROOR.PosOld = Ro.Position
					ROOR.PosNew = Rn.Position
					PSC.roles = true
					c = true
				}
				if c {
					ALLRC := PSC.RC
					ALLRC = append(ALLRC, ROOR)
					PSC.RC = ALLRC
				}
			} else {
				var ROOR = &PSRolechange{}
				if R.Mk {
					ROOR.ID = R.ID
					ROOR.Made = true
					PSC.roles = true
					ROOR.Mper = GetRole(R.ID, Gnew).Permissions
					ROOR.Existcrisis = true

					ALLRC := PSC.RC
					ALLRC = append(ALLRC, ROOR)
					PSC.RC = ALLRC
				} else if R.Del {
					ROOR.ID = R.ID
					ROOR.Deleted = true
					PSC.roles = true
					ROOR.Dperm = GetRole(R.ID, Gold).Permissions
					ROOR.Existcrisis = true

					ALLRC := PSC.RC
					ALLRC = append(ALLRC, ROOR)
					PSC.RC = ALLRC
				}
			}
		}
		for _, M := range Men.Members {
			if M.User.ID == Ps {
				if !M.ExistCrisis {
					NO := GetMember(Ps, Gold)
					NM := GetMember(Ps, Gnew)
					if M.User.Username {
						PSC.member = true
						PSC.MC.name = true
						PSC.MC.Name.NName = NM.User.Username
						PSC.MC.Name.OName = NO.User.Username
					}
					if M.Nick {
						PSC.member = true
						PSC.MC.nick = true
						PSC.MC.Nick.NNick = NM.Nick
						PSC.MC.Nick.OBick = NO.Nick
					}
					if Men.owner {
						PSC.owner = true
						PSC.OldOwner = Gold.OwnerID
						PSC.NewOwner = Gnew.OwnerID
						if Gold.OwnerID == Ps || Gnew.OwnerID == Ps {
							PSC.member = true
							PSC.MC.Owner = true
							if Gold.OwnerID == Ps && Gnew.OwnerID != Ps {
								PSC.MC.Isownernow = false
							} else if Gold.OwnerID != Ps && Gnew.OwnerID == Ps {
								PSC.MC.Isownernow = true
							}
						}
					}
				}
			}
		}
		if PSC.channels || PSC.member || PSC.roles || PSC.owner {
			ALLPSC = append(ALLPSC, PSC)
			is = true
		}
	}
	return ALLPSC, is
}

func HandleCase(pscmap []*PrimeSuspectChange, bichange *CompiledChange)
